//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2020 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package test

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/client/batch"
	"github.com/semi-technologies/weaviate/client/objects"
	"github.com/semi-technologies/weaviate/client/schema"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema/crossref"
	"github.com/semi-technologies/weaviate/entities/schema/kind"
	"github.com/semi-technologies/weaviate/test/acceptance/helper"
	"github.com/stretchr/testify/assert"
)

func Test_GraphQL(t *testing.T) {
	t.Run("setup test schema", addTestSchema)
	t.Run("import test data (city, country, airport)", addTestDataCityAirport)
	t.Run("import test data (companies)", addTestDataCompanies)
	t.Run("import test data (person)", addTestDataPersons)
	t.Run("import test data (custom vector class)", addTestDataCVC)

	// tests
	t.Run("getting objects", gettingObjects)
	t.Run("getting objects with filters", gettingObjectsWithFilters)
	t.Run("getting objects with geo filters", gettingObjectsWithGeoFilters)
	t.Run("getting objects with grouping", gettingObjectsWithGrouping)
	t.Run("getting objects with additional props", gettingObjectsWithAdditionalProps)

	// tear down
	deleteObjectClass(t, "Person")
	deleteObjectClass(t, "Country")
	deleteObjectClass(t, "City")
	deleteObjectClass(t, "Airport")
	deleteObjectClass(t, "Company")

	// only run after everything else is deleted, this way, we can also run an
	// all-class Explore since all vectors which are now left have the same
	// dimensions.
	t.Run("getting objects with custom vectors", gettingObjectsWithCustomVectors)
	deleteObjectClass(t, "CustomVectorClass")
}

func createObjectClass(t *testing.T, class *models.Class) {
	params := schema.NewSchemaObjectsCreateParams().WithObjectClass(class)
	resp, err := helper.Client(t).Schema.SchemaObjectsCreate(params, nil)
	helper.AssertRequestOk(t, resp, err, nil)
}

func createObject(t *testing.T, object *models.Object) {
	params := objects.NewObjectsCreateParams().WithBody(object)
	resp, err := helper.Client(t).Objects.ObjectsCreate(params, nil)
	helper.AssertRequestOk(t, resp, err, nil)
}

func createObjectsBatch(t *testing.T, objects []*models.Object) {
	params := batch.NewBatchObjectsCreateParams().
		WithBody(batch.BatchObjectsCreateBody{
			Objects: objects,
		})
	resp, err := helper.Client(t).Batch.BatchObjectsCreate(params, nil)
	helper.AssertRequestOk(t, resp, err, nil)
	for _, elem := range resp.Payload {
		assert.Nil(t, elem.Result.Errors)
	}
}

func deleteObjectClass(t *testing.T, class string) {
	delParams := schema.NewSchemaObjectsDeleteParams().WithClassName(class)
	delRes, err := helper.Client(t).Schema.SchemaObjectsDelete(delParams, nil)
	helper.AssertRequestOk(t, delRes, err, nil)
}

func addTestSchema(t *testing.T) {
	createObjectClass(t, &models.Class{
		Class: "Country",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"vectorizeClassName": true,
			},
		},
		Properties: []*models.Property{
			&models.Property{
				Name:     "name",
				DataType: []string{"string"},
			},
		},
	})

	createObjectClass(t, &models.Class{
		Class: "City",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"vectorizeClassName": true,
			},
		},
		Properties: []*models.Property{
			&models.Property{
				Name:     "name",
				DataType: []string{"string"},
			},
			&models.Property{
				Name:     "inCountry",
				DataType: []string{"Country"},
			},
			&models.Property{
				Name:     "population",
				DataType: []string{"int"},
			},
			&models.Property{
				Name:     "location",
				DataType: []string{"geoCoordinates"},
			},
		},
	})

	createObjectClass(t, &models.Class{
		Class: "Airport",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"vectorizeClassName": true,
			},
		},
		Properties: []*models.Property{
			&models.Property{
				Name:     "code",
				DataType: []string{"string"},
			},
			&models.Property{
				Name:     "phone",
				DataType: []string{"phoneNumber"},
			},
			&models.Property{
				Name:     "inCity",
				DataType: []string{"City"},
			},
		},
	})

	createObjectClass(t, &models.Class{
		Class: "Company",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"vectorizeClassName": false,
			},
		},
		Properties: []*models.Property{
			&models.Property{
				Name:     "name",
				DataType: []string{"string"},
				ModuleConfig: map[string]interface{}{
					"text2vec-contextionary": map[string]interface{}{
						"vectorizePropertyName": false,
					},
				},
			},
			&models.Property{
				Name:     "inCity",
				DataType: []string{"City"},
				ModuleConfig: map[string]interface{}{
					"text2vec-contextionary": map[string]interface{}{
						"vectorizePropertyName": false,
					},
				},
			},
		},
	})

	createObjectClass(t, &models.Class{
		Class: "Person",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"vectorizeClassName": false,
			},
		},
		Properties: []*models.Property{
			&models.Property{
				Name:     "name",
				DataType: []string{"string"},
				ModuleConfig: map[string]interface{}{
					"text2vec-contextionary": map[string]interface{}{
						"vectorizePropertyName": false,
					},
				},
			},
			&models.Property{
				Name:     "livesIn",
				DataType: []string{"City"},
				ModuleConfig: map[string]interface{}{
					"text2vec-contextionary": map[string]interface{}{
						"vectorizePropertyName": false,
					},
				},
			},
		},
	})

	createObjectClass(t, &models.Class{
		Class:      "CustomVectorClass",
		Vectorizer: "none",
		Properties: []*models.Property{},
	})
}

const (
	netherlands strfmt.UUID = "67b79643-cf8b-4b22-b206-6e63dbb4e57a"
	germany     strfmt.UUID = "561eea29-b733-4079-b50b-cfabd78190b7"
	amsterdam   strfmt.UUID = "8f5f8e44-d348-459c-88b1-c1a44bb8f8be"
	rotterdam   strfmt.UUID = "660db307-a163-41d2-8182-560782cd018f"
	berlin      strfmt.UUID = "9b9cbea5-e87e-4cd0-89af-e2f424fd52d6"
	dusseldorf  strfmt.UUID = "6ffb03f8-a853-4ec5-a5d8-302e45aaaf13"
	nullisland  strfmt.UUID = "823abeca-eef3-41c7-b587-7a6977b08003"
	airport1    strfmt.UUID = "4770bb19-20fd-406e-ac64-9dac54c27a0f"
	airport2    strfmt.UUID = "cad6ab9b-5bb9-4388-a933-a5bdfd23db37"
	airport3    strfmt.UUID = "55a4dbbb-e2af-4b2a-901d-98146d1eeca7"
	airport4    strfmt.UUID = "62d15920-b546-4844-bc87-3ae33268fab5"
	cvc1        strfmt.UUID = "1ffeb3e1-1258-4c2a-afc3-55543f6c44b8"
	cvc2        strfmt.UUID = "df22e5c4-5d17-49f9-a71d-f392a82bc086"
	cvc3        strfmt.UUID = "c28a039a-d509-4c2e-940a-8b109e5bebf4"
)

func addTestDataCityAirport(t *testing.T) {
	// countries
	createObject(t, &models.Object{
		Class: "Country",
		ID:    netherlands,
		Properties: map[string]interface{}{
			"name": "Netherlands",
		},
	})
	createObject(t, &models.Object{
		Class: "Country",
		ID:    germany,
		Properties: map[string]interface{}{
			"name": "Germany",
		},
	})

	// cities
	createObject(t, &models.Object{
		Class: "City",
		ID:    amsterdam,
		Properties: map[string]interface{}{
			"name":       "Amsterdam",
			"population": 1800000,
			"location": map[string]interface{}{
				"latitude":  52.366667,
				"longitude": 4.9,
			},
			"inCountry": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", netherlands, kind.Object).String(),
				},
			},
		},
	})
	createObject(t, &models.Object{
		Class: "City",
		ID:    rotterdam,
		Properties: map[string]interface{}{
			"name":       "Rotterdam",
			"population": 600000,
			"inCountry": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", netherlands, kind.Object).String(),
				},
			},
		},
	})
	createObject(t, &models.Object{
		Class: "City",
		ID:    berlin,
		Properties: map[string]interface{}{
			"name":       "Berlin",
			"population": 3470000,
			"inCountry": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", germany, kind.Object).String(),
				},
			},
		},
	})
	createObject(t, &models.Object{
		Class: "City",
		ID:    dusseldorf,
		Properties: map[string]interface{}{
			"name":       "Dusseldorf",
			"population": 600000,
			"inCountry": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", germany, kind.Object).String(),
				},
			},
			"location": map[string]interface{}{
				"latitude":  51.225556,
				"longitude": 6.782778,
			},
		},
	})

	createObject(t, &models.Object{
		Class: "City",
		ID:    nullisland,
		Properties: map[string]interface{}{
			"name":       "Null Island",
			"population": 0,
			"location": map[string]interface{}{
				"latitude":  0,
				"longitude": 0,
			},
		},
	})

	// airports
	createObject(t, &models.Object{
		Class: "Airport",
		ID:    airport1,
		Properties: map[string]interface{}{
			"code": "10000",
			"phone": map[string]interface{}{
				"input": "+311234567",
			},
			"inCity": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", amsterdam, kind.Object).String(),
				},
			},
		},
	})
	createObject(t, &models.Object{
		Class: "Airport",
		ID:    airport2,
		Properties: map[string]interface{}{
			"code": "20000",
			"inCity": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", rotterdam, kind.Object).String(),
				},
			},
		},
	})
	createObject(t, &models.Object{
		Class: "Airport",
		ID:    airport3,
		Properties: map[string]interface{}{
			"code": "30000",
			"inCity": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", dusseldorf, kind.Object).String(),
				},
			},
		},
	})
	createObject(t, &models.Object{
		Class: "Airport",
		ID:    airport4,
		Properties: map[string]interface{}{
			"code": "40000",
			"inCity": []interface{}{
				map[string]interface{}{
					"beacon": crossref.New("localhost", berlin, kind.Object).String(),
				},
			},
		},
	})

	// wait for consistency
	assertGetObjectEventually(t, airport1)
	assertGetObjectEventually(t, airport2)
	assertGetObjectEventually(t, airport3)
	assertGetObjectEventually(t, airport4)

	// give cache some time to become hot
	time.Sleep(2 * time.Second)
}

func addTestDataCompanies(t *testing.T) {
	var (
		microsoft1 strfmt.UUID = "cfa3b21e-ca4f-4db7-a432-7fc6a23c534d"
		microsoft2 strfmt.UUID = "8f75ed97-39dd-4294-bff7-ecabd7923062"
		microsoft3 strfmt.UUID = "f343f51d-7e05-4084-bd66-d504db3b6bec"
		apple1     strfmt.UUID = "477fec91-1292-4928-8f53-f0ff49c76900"
		apple2     strfmt.UUID = "bb2cfdba-d4ba-4cf8-abda-e719ef35ac33"
		apple3     strfmt.UUID = "b71d2b4c-3da1-4684-9c5e-aabd2a4f2998"
		google1    strfmt.UUID = "8c2e21fc-46fe-4999-b41c-a800595129af"
		google2    strfmt.UUID = "62b969c6-f184-4be0-8c40-7470af417cfc"
		google3    strfmt.UUID = "c7829929-2037-4420-acbc-a433269feb93"
	)

	type companyTemplate struct {
		id     strfmt.UUID
		name   string
		inCity []strfmt.UUID
	}

	companies := []companyTemplate{
		companyTemplate{id: microsoft1, name: "Microsoft Inc.", inCity: []strfmt.UUID{dusseldorf}},
		companyTemplate{id: microsoft2, name: "Microsoft Incorporated", inCity: []strfmt.UUID{dusseldorf, amsterdam}},
		companyTemplate{id: microsoft3, name: "Microsoft", inCity: []strfmt.UUID{berlin}},
		companyTemplate{id: apple1, name: "Apple Inc."},
		companyTemplate{id: apple2, name: "Apple Incorporated"},
		companyTemplate{id: apple3, name: "Apple"},
		companyTemplate{id: google1, name: "Google Inc."},
		companyTemplate{id: google2, name: "Google Incorporated"},
		companyTemplate{id: google3, name: "Google"},
	}

	// companies
	for _, company := range companies {
		inCity := []interface{}{}
		for _, c := range company.inCity {
			inCity = append(inCity,
				map[string]interface{}{
					"beacon": crossref.New("localhost", c, kind.Object).String(),
				})
		}

		createObject(t, &models.Object{
			Class: "Company",
			ID:    company.id,
			Properties: map[string]interface{}{
				"inCity": inCity,
				"name":   company.name,
			},
		})
	}

	assertGetObjectEventually(t, companies[len(companies)-1].id)
}

func addTestDataPersons(t *testing.T) {
	var (
		alice strfmt.UUID = "5d0fa6ee-21c4-4b46-a735-f0208717837d"
		bob   strfmt.UUID = "8615585a-2960-482d-b19d-8bee98ade52c"
		john  strfmt.UUID = "3ef44474-b5e5-455d-91dc-d917b5b76165"
		petra strfmt.UUID = "15d222c9-8c36-464b-bedb-113faa1c1e4c"
	)

	type personTemplate struct {
		id      strfmt.UUID
		name    string
		livesIn []strfmt.UUID
	}

	companies := []personTemplate{
		personTemplate{id: alice, name: "Alice", livesIn: []strfmt.UUID{}},
		personTemplate{id: bob, name: "Bob", livesIn: []strfmt.UUID{amsterdam}},
		personTemplate{id: john, name: "John", livesIn: []strfmt.UUID{amsterdam, berlin}},
		personTemplate{id: petra, name: "Petra", livesIn: []strfmt.UUID{amsterdam, berlin, dusseldorf}},
	}

	// companies
	for _, person := range companies {
		livesIn := []interface{}{}
		for _, c := range person.livesIn {
			livesIn = append(livesIn,
				map[string]interface{}{
					"beacon": crossref.New("localhost", c, kind.Object).String(),
				})
		}

		createObject(t, &models.Object{
			Class: "Person",
			ID:    person.id,
			Properties: map[string]interface{}{
				"livesIn": livesIn,
				"name":    person.name,
			},
		})
	}

	assertGetObjectEventually(t, companies[len(companies)-1].id)
}

func addTestDataCVC(t *testing.T) {
	// add one object indivdually
	createObject(t, &models.Object{
		Class:  "CustomVectorClass",
		ID:     cvc1,
		Vector: []float32{1.1, 1.1, 1.1},
	})

	assertGetObjectEventually(t, cvc1)

	createObjectsBatch(t, []*models.Object{
		{
			Class:  "CustomVectorClass",
			ID:     cvc2,
			Vector: []float32{1.1, 1.1, 0.1},
		},
		{
			Class:  "CustomVectorClass",
			ID:     cvc3,
			Vector: []float32{1.1, 0, 0},
		},
	})
	assertGetObjectEventually(t, cvc3)
}
