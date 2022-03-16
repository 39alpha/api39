package survey

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
)

func OpenDatabase(ctx iris.Context) {
	settings := sqlite.ConnectionURL{
		Database: `surveys.db`,
	}

	sess, err := sqlite.Open(settings)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "cannot open survey database",
		})
		return
	}

	ctx.Values().Set("db_session", sess)
	ctx.Next()
}

func ListSurveys(ctx iris.Context) {
	sess, ok := ctx.Values().Get("db_session").(db.Session)
	if !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no database loaded",
		})
		return
	}

	var ids []iris.Map
	if err := sess.SQL().Select("id").From("survey").Iterator().All(&ids); err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "failed to get surveys",
		})
		return
	}

	_, _ = ctx.JSON(iris.Map{"surveys": ids})
}

func GetSurvey(ctx iris.Context) {
	survey_id := ctx.Params().Get("id")
	sess, ok := ctx.Values().Get("db_session").(db.Session)
	if !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no database loaded",
		})
		return
	}

	iter := sess.SQL().
		Select(
			"s.id as survey_id",
			"s.title as title",
			"s.description as description",
			"q.id as question_id",
			"q.type as type",
			"q.statement as statement",
			"a.id as answer_id",
			"a.value as value",
		).
		From("survey as s").
		LeftJoin("question as q").
		On("s.id = q.survey_id").
		LeftJoin("answer as a").
		On("q.id = a.question_id").
		Where("s.id = ?", survey_id).
		OrderBy("s.id", "q.id", "a.id").
		Iterator()

	var results []map[string]interface{}
	err := iter.All(&results)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "failed to fetch survey",
		})
		return
	}
	if len(results) == 0 {
		ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
			"message": fmt.Sprintf("no survey found with id '%s'", survey_id),
		})
		return
	}

	survey := iris.Map{
		"id":          results[0]["survey_id"],
		"title":       results[0]["title"],
		"description": results[0]["description"],
		"questions":   make([]iris.Map, 0),
	}

	if results[0]["question_id"] != nil {
		question := iris.Map{}
		for _, row := range results {
			question_id := row["question_id"]
			if question["id"] == nil {
				question = iris.Map{
					"id":        question_id,
					"type":      row["type"],
					"statement": row["statement"],
					"answers":   make([]iris.Map, 0),
				}
			} else if question["id"] != question_id {
				survey["questions"] = append(
					survey["questions"].([]iris.Map),
					question,
				)
				question = iris.Map{
					"id":        question_id,
					"type":      row["type"],
					"statement": row["statement"],
					"answers":   make([]iris.Map, 0),
				}
			}

			if row["value"] != nil {
				question["answers"] = append(
					question["answers"].([]iris.Map),
					iris.Map{
						"id": row["answer_id"],
						"value": row["value"].(string),
					},
				)
			}
		}

		if len(question) != 0 {
			survey["questions"] = append(
				survey["questions"].([]iris.Map),
				question,
			)
		}
	}

	_, _ = ctx.JSON(survey)
}

func AddSurveyResponses(ctx iris.Context) {
	body, ok := ctx.Values().Get("JSONBody").(iris.Map)
	if !ok {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("Bad request body"))
		return
	}

	survey_id := ctx.Params().Get("id")
	sess, ok := ctx.Values().Get("db_session").(db.Session)
	if !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no database loaded",
		})
		return
	}

	err := sess.Tx(func (sess db.Session) error {
		submission_id, err := addSubmission(sess, survey_id)
		if err != nil {
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"error":   "failed to get submission id",
				"message": fmt.Sprintf("%v", err),
			})
			return err
		}
	
		if content, ok := body["responses"]; !ok {
			ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
				"message": "ill-formed response",
			})
			return fmt.Errorf("ill-formed response")
		} else if responses, ok := content.([]interface {}); !ok {
			ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
				"message": "ill-formed 'responses' field",
			})
			return fmt.Errorf("ill-formed 'responses' field")
		} else {
			inserter := sess.SQL().
				InsertInto("response").
				Columns("submission_id", "question_id", "response").
				Batch(len(responses))
	
			for _, response := range responses {
				if entry, ok := response.(iris.Map); !ok {
					ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
						"error": "ill-formed survey response",
					})
					return fmt.Errorf("ill-formed survey response")
				} else if question_id, ok := entry["id"].(float64); !ok {
					ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
						"error": "ill-formed 'id' field in survey response",
					})
					return fmt.Errorf("ill-formed 'id' field in survey response")
				} else if value, ok := entry["response"].(string); !ok {
					ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
						"error": "ill-formed 'response' field in survey response",
					})
					return fmt.Errorf("ill-formed 'response' field in survey response")
				} else {
					inserter = inserter.Values(submission_id, question_id, value)
				}
			}
	
			inserter.Done()
	
			if err := inserter.Wait(); err != nil {
				ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
					"error":   "failed to insert the responses",
					"message": fmt.Sprintf("%v", err),
				})
				return err
			}
		}

		return nil
	});

	if err != nil {
		return
	}

	_, _ = ctx.JSON(iris.Map{"message": "success"})
}

func addSubmission(sess db.Session, survey_id string) (int, error) {
	rows, err := sess.SQL().
		InsertInto("submission").
		Columns("survey_id").
		Values(survey_id).
		Returning("id").
		Query()

	if err != nil {
		return 0, nil
	}
	defer rows.Close()

	var submission_id int
	rows.Next()
	if err := rows.Scan(&submission_id); err != nil {
		return 0, err
	}

	return submission_id, nil
}
