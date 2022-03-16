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
			"a.answer as answer",
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
					"answers":   make([]string, 0),
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
					"answers":   make([]string, 0),
				}
			}

			if row["answer"] != nil {
				question["answers"] = append(
					question["answers"].([]string),
					row["answer"].(string),
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

func PutSurveyResponses(ctx iris.Context) {
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

	submission_id, err := addSubmission(sess, survey_id)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error":   "failed to get submission id",
			"message": fmt.Sprintf("%v", err),
		})
	}

	if responses, ok := body["responses"]; !ok {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": "ill-formed response",
		})
		return
	} else {
		inserter := sess.SQL().
			InsertInto("response").
			Columns("submission_id", "question_id", "answer").
			Batch(len(body))

		for _, response := range responses.([]interface{}) {
			entry := response.(iris.Map)
			question_id := int(entry["question_id"].(float64))
			answer := entry["answer"].(string)
			inserter = inserter.Values(submission_id, question_id, answer)
		}

		inserter.Done()

		if err := inserter.Wait(); err != nil {
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"error":   "failed to insert the responses",
				"message": fmt.Sprintf("%v", err),
			})
			return
		}
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
