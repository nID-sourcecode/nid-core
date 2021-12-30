package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// DefaultModelUserFirstAndLastName is a default representation of the model
func (m *AccessModelDB) DefaultModelUserFirstAndLastName(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		AudienceID:  audienceDB.DefaultModelPrimary(namespace).ID,
		Description: "Voor- en achternaam",
		Hash:        "3bf8826ff2805c4822536e609ae9d1dcb3c3b9b89e077f0d73963c60c1a408e9",
		Type:        AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			JSONModel: `
			{
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"f": ["firstName", "lastName"],
					"p": {
						"filter": {
							"pseudonym": {
								"eq": "$$nid:subject$$"
							}
						}
					}
				}
			}
		`,
			Path: "/gql",
		},
		Name: "databron:name",
	}
}

// DefaultModelOptionalBankAccounts is an optional default representation of the model
func (m *AccessModelDB) DefaultModelOptionalBankAccounts(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		AudienceID:  audienceDB.DefaultModelPrimary(namespace).ID,
		Description: "Bankrekeningen",
		Hash:        "791a985542b63d6593830bda69258da190d79cfdd28482cc5e6f85031c375189",
		Type:        AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			JSONModel: `
			{
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"m": {
						"bankAccounts": "#B"
					},
					"f": [],
					"p": {
						"filter": {
							"pseudonym": {
								"eq": "$$nid:subject$$"
							}
						}
					}
				},
				"B": {
					"m": {
						"savingsAccounts": "#S"
					},
					"f": [
						"accountNumber",
						"amount"
					] 
				},
				"S": {
					"f": [
						"amount",
						"name"
					]
				}
			}		
			`,
			Path: "/gql",
		},
		Name: "databron:bankaccounts",
	}
}

// DefaultModelUserAddresses is a default representation of an access model
func (m *AccessModelDB) DefaultModelUserAddresses(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		AudienceID:  audienceDB.DefaultModelSecondary(namespace).ID,
		Description: "Adresgegevens",
		Hash:        "3bf8826ff2805c4822536e609ae9d1dcb3c3b9b89e077f0d73963c60c1a408e9",
		Type:        AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			JSONModel: `
			{
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"m": {
						"contactDetails": "#C"
					},
					"f": [],
					"p": {
						"filter": {
							"pseudonym": {
								"eq": "$$nid:subject$$"
							}
						}
					}
				},
				"C": {
					"f": [
						"phone"
					] 
				}
			}
		`,
			Path: "/gql",
		},
		Name: "databron:addresses",
	}
}

// DefaultModelOptionalUserAddressContactDetails is a default representation of an access model
func (m *AccessModelDB) DefaultModelOptionalUserAddressContactDetails(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		AudienceID:  audienceDB.DefaultModelSecondary(namespace).ID,
		Description: "Contactgegevens",
		Hash:        "791a985542b63d6593830bda69258da190d79cfdd28482cc5e6f85031c375189",
		Type:        AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			JSONModel: `
				{
					"r": {
						"m": {
							"users": "#U"
						}
					},
					"U": {
						"m": {
							"contactDetails": "#C"
						},
						"f": [],
						"p": {
							"filter": {
								"pseudonym": {
									"eq": "$$nid:subject$$"
								}
							}
						}
					},
					"C": {
						"m": {
							"address": "#A"
						},
						"f": [] 
					},
					"A": {
						"f": [
							"houseNumber"
						]
					}
				}
			`,
			Path: "/gql",
		},
		Name: "databron:contactdetails",
	}
}

// DefaultModelUserFirstAndLastNameByBSN is an access model that represents getting a name by BSN
func (m *AccessModelDB) DefaultModelUserFirstAndLastNameByBSN(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		AudienceID:  audienceDB.DefaultModelPrimary(namespace).ID,
		Description: "Voor- en achternaam middels BSN",
		Hash:        "3bf8826ff2805c4822536e609ae9d1dcb3c3b9b89e077f0d73963c60c1a408e9",
		Type:        AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			JSONModel: `
			{
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"f": ["firstName", "lastName"],
					"p": {
						"filter": {
							"bsn": {
								"eq": "$$nid:bsn$$"
							}
						}
					}
				}
			}
		`,
			Path: "/gql",
		},
		Name: "databron:name:bsn",
	}
}

// DefaultModelOptionalBankAccountsByBSN is an access model that represents getting a bank account by BSN
func (m *AccessModelDB) DefaultModelOptionalBankAccountsByBSN(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		AudienceID:  audienceDB.DefaultModelPrimary(namespace).ID,
		Description: "Bankrekeningen middels BSN",
		Hash:        "791a985542b63d6593830bda69258da190d79cfdd28482cc5e6f85031c375189",
		Type:        AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			JSONModel: `
			{
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"m": {
						"bankAccounts": "#B"
					},
					"f": [],
					"p": {
						"filter": {
							"bsn": {
								"eq": "$$nid:bsn$$"
							}
						}
					}
				},
				"B": {
					"m": {
						"savingsAccounts": "#S"
					},
					"f": [
						"accountNumber",
						"amount"
					] 
				},
				"S": {
					"f": [
						"amount",
						"name"
					]
				}
			}		
			`,
			Path: "/gql",
		},
		Name: "databron:bankaccounts:bsn",
	}
}

// DefaultModelContactable is the default model for the contactable endpoint on the information service
func (m *AccessModelDB) DefaultModelContactable(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		ID:          uuid.FromStringOrNil("4d0888c7-56df-4369-9395-00692518babf"),
		AudienceID:  audienceDB.DefaultModelInformationService(namespace).ID,
		Description: "Heeft gebruiker contactgegevens (ja/nee)",
		Hash:        "884240858a33744a106fdc13edf463703ee50f0629c3573596d6617cb39fbf2f",
		Name:        "information:contactable",
		RestAccessModel: &RestAccessModel{
			Body:   "",
			Method: "GET",
			Path:   "/v1/info/contact-details/contactable",
			Query:  "{}",
		},
		Type: AccessModelTypeREST,
	}
}

// DefaultModelHasAddress is the default model for the has-address endpoint on the information service
func (m *AccessModelDB) DefaultModelHasAddress(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		ID:          uuid.FromStringOrNil("fae3d71e-4121-49ff-8fd3-a66eac4bc4d3"),
		AudienceID:  audienceDB.DefaultModelInformationService(namespace).ID,
		Description: "Heeft gebruiker een postadres (ja/nee)",
		Hash:        "51ce4803bb996acbe69c08aed9d445adddf77f934d744a983faec293f8c39f31",
		Name:        "information:hasaddress",
		RestAccessModel: &RestAccessModel{
			Body:   "",
			Method: "GET",
			Path:   "/v1/info/contact-details/has-address",
			Query:  "{}",
		},
		Type: AccessModelTypeREST,
	}
}

// DefaultModelHasPositiveBankAccountBalance is the default model for the positive balance endpoint on the information service
func (m *AccessModelDB) DefaultModelHasPositiveBankAccountBalance(namespace string) AccessModel {
	audienceDB := &AudienceDB{}

	return AccessModel{
		ID:          uuid.FromStringOrNil("6f611ed8-36ef-4ded-92cb-d5577f9538bf"),
		AudienceID:  audienceDB.DefaultModelInformationService(namespace).ID,
		Description: "Heeft gebruiker een positief saldo (ja/nee)",
		Hash:        "aaecd8f3facce35953f0b9f416e1fdeb1f114210480e98c7df814826189dccaf",
		Name:        "information:bankaccountpositive",
		RestAccessModel: &RestAccessModel{
			Body:   "",
			Method: "GET",
			Path:   "/v1/info/bank-account/positive",
			Query:  "{}",
		},
		Type: AccessModelTypeREST,
	}
}

// CreateAccessModel inserts the accessmodel
func (m *AccessModelDB) CreateAccessModel(model *AccessModel) error {
	err := m.Db.Create(&model).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAccessModelsByIDs retrieves the accessmodels from the given ids
func (m *AccessModelDB) GetAccessModelsByIDs(suppliedAccessModels []string) ([]*AccessModel, error) {
	var accessModels []*AccessModel

	err := m.Db.Where(`id IN (?)`, suppliedAccessModels).Find(&accessModels).Error
	if err != nil {
		return nil, err
	}

	// check if accessmodels are found. gorm does not return this error when filling a slice
	if len(accessModels) != len(suppliedAccessModels) {
		return nil, gorm.ErrRecordNotFound
	}

	return accessModels, nil
}

// GetAccessModelByAudienceWithScope retrieves an accessmodel given a state and audience
func (m *AccessModelDB) GetAccessModelByAudienceWithScope(name string, hash string, audience *Audience) (*AccessModel, error) {
	var accessModel AccessModel
	err := m.Db.Find(&accessModel, "name = ? AND hash = ? AND audience_id = ?", name, hash, audience.ID).Error
	if err != nil {
		return nil, err
	}

	return &accessModel, nil
}

// GetAccessModelsByAudience retrieves the accessmodels of an audience
func (m *AccessModelDB) GetAccessModelsByAudience(preloadGqlAndRest bool, audience *Audience) ([]*AccessModel, error) {
	var accessModels []*AccessModel
	var err error
	if preloadGqlAndRest {
		err = m.Db.Preload("GqlAccessModel").Preload("RestAccessModel").Find(&accessModels, "audience_id = ?", audience.ID).Error
	} else {
		err = m.Db.Find(&accessModels, "audience_id = ?", audience.ID).Error
	}

	if err != nil {
		return nil, err
	}

	// check if accessmodels are found. gorm does not return this error when filling a slice
	if len(accessModels) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return accessModels, nil
}
