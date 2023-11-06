// Package internal contains the internal logic for luarunner
package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	authModels "github.com/nID-sourcecode/nid-core/svc/auth/models"
	auth "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
	"github.com/nID-sourcecode/nid-core/svc/luarunner/internal/runner"
	"github.com/nID-sourcecode/nid-core/svc/luarunner/models"
)

// LuaRunnerService handles the requests for the lua scripts.
type LuaRunnerService struct {
	scriptRunner runner.ScriptRunner

	luaRunnerDB      *LuaRunnerDB
	authDB           *gorm.DB
	authClientDB     *authModels.ClientDB
	redirectTargetDB *authModels.RedirectTargetDB
	audienceDB       *authModels.AudienceDB
}

// NewLuaRunnerService returns a new instance of LuaRunnerService
func NewLuaRunnerService(authDB *gorm.DB, authClient *auth.AuthClient, db *LuaRunnerDB) *LuaRunnerService {
	luaRunner := runner.NewLuaRunner(authClient)

	return &LuaRunnerService{
		scriptRunner:     &luaRunner,
		luaRunnerDB:      db,
		authDB:           authDB,
		authClientDB:     authModels.NewClientDB(authDB),
		redirectTargetDB: authModels.NewRedirectTargetDB(authDB),
		audienceDB:       authModels.NewAudienceDB(authDB),
	}
}

// ScriptRunner struct contains information for the script that will run.
type ScriptRunner struct {
	organisation   *models.Organisation
	audience       *authModels.Audience
	redirectTarget *authModels.RedirectTarget
	script         string
}

type callbackRequest struct {
	OrganisatieID     string `json:"organisatieId"`
	OrganisatieIDType string `json:"organisatieIdType"`
	Timestamp         string `json:"timestamp"`
	AbonnementID      string `json:"abonnementId"`
	EventType         string `json:"eventType"`
	RecordID          string `json:"recordId"`
}

var (
	validateRecordID    = regexp.MustCompile(`^[\w+:/?\d\-.]+$`).MatchString  //nolint:gochecknoglobals
	validateRunJSONBody = regexp.MustCompile(`^[\p{L}\d\s\-/]+$`).MatchString //nolint:gochecknoglobals
)

var (
	errRecordIDIsInvalid    = fmt.Errorf("recordId is invalid")
	errOrganisationNotFound = fmt.Errorf("could not find organisation from organisatieId")
)

// ErrJSONBodyIsInvalid error message for when the json body is invalid for laurunner.
var ErrJSONBodyIsInvalid = fmt.Errorf("json body is invalid: allowed characters: alpha numeric, forward slash, hyphen and spaces")

// HTTPCallback endpoint that handles the requests for callback method of luarunner service.
func (l *LuaRunnerService) HTTPCallback(w http.ResponseWriter, r *http.Request) {
	if !handleMethod(w, r, "POST") {
		return
	}

	var callbackRequest callbackRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&callbackRequest); err != nil {
		log.WithError(err).Error("decoding callback request")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("500: internal server error"))
		if err != nil {
			log.WithError(err).Error("writing")
		}
		return
	}

	err := l.Callback(r.Context(), &callbackRequest)
	if err != nil {
		callbackRespondWithError(w, err)
		return
	}

	// returning 202 instead of 200 to avoid confusion with the caller of callback.
	w.WriteHeader(http.StatusAccepted)
}

// Callback is used to run lua scripts which are assigned to the given organisation.
func (l *LuaRunnerService) Callback(ctx context.Context, callbackRequest *callbackRequest) error {
	// user input check
	if !validateRecordID(callbackRequest.RecordID) {
		log.Info(fmt.Sprintf("recordId: %s is not valid", callbackRequest.RecordID))
		return errRecordIDIsInvalid
	}

	input := make(map[string]interface{})

	organisation, err := l.luaRunnerDB.OrganisationDB.GetOrganisationWithUzoviName(callbackRequest.OrganisatieID)
	if err != nil {
		log.WithError(err).Error("getting organisation with uzovi name")
		return errOrganisationNotFound
	}

	redirectTarget, err := l.redirectTargetDB.Get(ctx, organisation.RedirectTargetID)
	if err != nil {
		log.WithError(err).Error("getting redirect urls from database")
		return err
	}

	audience, err := l.audienceDB.Get(ctx, organisation.AudienceID)
	if err != nil {
		log.WithError(err).Error("getting audience from database")
		return err
	}

	input["clientID"] = redirectTarget.ClientID.String()
	input["redirectID"] = redirectTarget.RedirectTarget
	input["audience"] = audience.Audience

	recordSplit := strings.Split(callbackRequest.RecordID, "/")
	input["wlzindicatieID"] = recordSplit[len(recordSplit)-1]

	for _, s := range organisation.Scripts {
		if s.Script == "" {
			log.WithField("ScriptID", s.ID).Error("given script from database was empty")
			continue
		}

		err = l.scriptRunner.RunScript(ctx, s.Script, input)
		if err != nil {
			log.WithError(err).Error("running lua script")
			return err
		}
	}

	return nil
}

// Run runs the given script of the given organisation
func (l *LuaRunnerService) Run(c *gin.Context) {
	if !handleMethod(c.Writer, c.Request, "POST") {
		return
	}
	organisationID := c.Param("organisation")
	scriptID := c.Param("script")

	ctx := c.Request.Context()

	if organisationID == "" || scriptID == "" {
		writeError(
			c, http.StatusNotFound,
			fmt.Sprintf("Script \"%s\" for organisation \"%s\" not found", scriptID, organisationID),
		)
		return
	}

	organisation, err := l.luaRunnerDB.OrganisationDB.Get(ctx, uuid.FromStringOrNil(organisationID))
	if err != nil {
		log.WithError(err).WithField("organisationID", organisationID).Error("getting organisation with id")
		writeError(c, http.StatusNotFound, fmt.Sprintf("Organisation \"%s\" not found", organisationID))
		return
	}

	scriptRunner, err := l.getScriptRunner(ctx, organisation, scriptID)
	if err != nil {
		log.WithError(err).Error("getting organisation details")
		writeError(
			c, http.StatusNotFound,
			fmt.Sprintf("Organisation details for organisation \"%s\" not found", organisationID),
		)
		return
	}

	var body bytes.Buffer
	copyBody := io.TeeReader(c.Request.Body, &body)

	err = l.ValidateJSONBody(copyBody)
	if err != nil {
		log.WithError(err).Error("validating json body")
		writeError(c, http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		return
	}

	input := make(map[string]interface{})

	// parse json body
	if err := json.NewDecoder(&body).Decode(&input); err != nil {
		log.WithError(err).Error("decoding body")
		writeError(c, http.StatusBadRequest, "Unable to parse body")
		return
	}

	input["clientID"] = scriptRunner.redirectTarget.ClientID.String()
	input["redirectID"] = scriptRunner.redirectTarget.RedirectTarget
	input["audience"] = scriptRunner.audience.Audience
	input["payload"] = scriptRunner.audience.Audience

	err = l.scriptRunner.RunScript(ctx, scriptRunner.script, input)

	if err != nil {
		log.WithError(err).Error("running lua script")
		writeError(c, http.StatusInternalServerError, "Failed running script")

		return
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}

func (l *LuaRunnerService) getScriptRunner(
	ctx context.Context, organisation *models.Organisation, scriptID string,
) (*ScriptRunner, error) {
	redirectTarget, err := l.redirectTargetDB.Get(ctx, organisation.RedirectTargetID)
	if err != nil {
		log.WithError(err).Error("getting redirect urls from database")
		return nil, err
	}

	audience, err := l.audienceDB.Get(ctx, organisation.AudienceID)
	if err != nil {
		log.WithError(err).Error("getting audience from database")
		return nil, err
	}

	script, err := l.luaRunnerDB.ScriptDB.Get(ctx, uuid.FromStringOrNil(scriptID))
	if err != nil {
		log.WithError(err).Error("getting script from database")
		return nil, err
	}

	scriptRunner := &ScriptRunner{
		organisation:   organisation,
		audience:       audience,
		redirectTarget: redirectTarget,
		script:         script.Script,
	}

	return scriptRunner, nil
}

// ValidateJSONBody validates if the json body meets the criteria of Luarunner
func (l *LuaRunnerService) ValidateJSONBody(body io.Reader) error {
	validateErr := ErrJSONBodyIsInvalid
	ok := true

	dec := json.NewDecoder(body)
	for {
		t, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if _, ok := t.(json.Delim); ok || t == nil {
			continue
		}

		s := fmt.Sprintf("%v", t)
		if !validateRunJSONBody(s) {
			validateErr = errors.Wrap(validateErr, fmt.Sprintf("%s is not valid", s))
			ok = false
		}
	}

	if !ok {
		return validateErr
	}

	return nil
}

func callbackRespondWithError(w http.ResponseWriter, err error) {
	log.WithError(err).Error("running callback method")

	if errors.Is(err, errRecordIDIsInvalid) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("recordID is invalid"))
		if err != nil {
			log.WithError(err).Error("writing http response")
		}
		return
	}

	if errors.Is(err, errOrganisationNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Bad Request For: organisatieId"))
		if err != nil {
			log.WithError(err).Error("writing http response")
		}
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte("500: internal server error"))
	if err != nil {
		log.WithError(err).Error("writing")
	}
}

func writeError(c *gin.Context, statusCode int, message string) {
	c.Writer.WriteHeader(statusCode)
	_, err := c.Writer.WriteString(message)
	if err != nil {
		log.WithError(err).Error("writing")
	}
}

func handleMethod(w http.ResponseWriter, r *http.Request, methodName string) bool {
	if r.Method != methodName {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("405: Method not allowed"))
		if err != nil {
			log.WithError(err).Error("tried to write a response")
		}
		return false
	}

	return true
}
