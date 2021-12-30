package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"

	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	authModels "lab.weave.nl/nid/nid-core/svc/auth/models"
	auth "lab.weave.nl/nid/nid-core/svc/auth/proto"
	"lab.weave.nl/nid/nid-core/svc/luarunner/runner"
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
func NewLuaRunnerService(authDB *gorm.DB, authClient auth.AuthClient, db *LuaRunnerDB) *LuaRunnerService {
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

type callbackRequest struct {
	OrganisatieID     string `json:"organisatieId"`
	OrganisatieIDType string `json:"organisatieIdType"`
	Timestamp         string `json:"timestamp"`
	AbonnementID      string `json:"abonnementId"`
	EventType         string `json:"eventType"`
	RecordID          string `json:"recordId"`
}

// HTTPCallback endpoint that handles the requests for callback method of luarunner service.
func (l *LuaRunnerService) HTTPCallback(w http.ResponseWriter, r *http.Request) {
	if !handleMethod(w, r, "POST") {
		return
	}

	var callbackRequest callbackRequest

	if err := json.NewDecoder(r.Body).Decode(&callbackRequest); err != nil {
		log.WithError(err).Error("decoding callback request")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("500: internal server error"))
		if err != nil {
			log.WithError(err).Error("writing")
		}
		return
	}

	if err := l.Callback(r.Context(), &callbackRequest); err != nil {
		log.WithError(err).Error("running callback method")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("500: internal server error"))
		if err != nil {
			log.WithError(err).Error("writing")
		}
		return
	}

	// returning 202 instead of 200 to avoid confusion with the caller of callback.
	w.WriteHeader(http.StatusAccepted)
}

// Callback is used to run lua scripts which are assigned to the given organisation.
func (l *LuaRunnerService) Callback(ctx context.Context, callbackRequest *callbackRequest) error {
	input := make(map[string]interface{})

	organisation, err := l.luaRunnerDB.OrganisationDB.GetOrganisationWithUzoviName(callbackRequest.OrganisatieID)
	if err != nil {
		log.WithError(err).Error("getting organisation with uzovi name")
		return err
	}

	l.redirectTargetDB.DB().(*gorm.DB).Where("")
	if err != nil {
		log.WithError(err).Error("getting client from database")
		return err
	}

	redirectTarget, err := l.redirectTargetDB.Get(ctx, organisation.RedirectTargetID)
	if err != nil {
		log.WithError(err).Error("getting redirect urls from database")
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

func handleMethod(w http.ResponseWriter, r *http.Request, methodName string) bool {
	if r.Method != methodName {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("405: Method not allowed"))
		if err != nil {
			log.WithError(err).Fatal("writing")
		}
		return false
	}

	return true
}
