package survey

import (
	"github.com/kataras/iris/v12"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
)

type Survey struct {
	Id string `db:"id,omitempty"`
	Title string `db:"title"`
	Description string `db:"description"`
}

type Question struct {
	Id int `db:"id,omitempty"`
	Type string `db:type`
	Statement string `db:"statement"`
		SurveyId string `db:"survey_id"`
}

type Answer struct {
	Id int `db:"id,omitempty"`
	Answer string `db:"answer"`
		QuestionId int `db:"question_id"`
}

type Response struct {
	Id int `db:"id,omitempty"`
	SurveyId string `db:"survey_id"`
	QuestionId string `db:"question_id"`
	Answer string `db:"answer"`
}

func OpenDatabase(ctx iris.Context) {
	settings := sqlite.ConnectionURL{
		Database: `surveys.db`,
	}

	db, err := sqlite.Open(settings)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "cannot open survey database",
		})
		return
	}

	ctx.Values().Set("surveydb", db)
	ctx.Next()
}

func ListSurveys(ctx iris.Context) {
	sess, ok := ctx.Values().Get("surveydb").(db.Session)
	if !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no database loaded",
		})
		return
	}

	surveyCollection := sess.Collection("survey")
	res := surveyCollection.Find()
	var surveys []Survey
	if err := res.All(&surveys); err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "failed to get surveys",
		})
		return
	}

	var ids []iris.Map
	for _, survey := range surveys {
		ids = append(ids, iris.Map{ "id": survey.Id })
	}
	_, _ = ctx.JSON(iris.Map{"surveys": ids})
}

func GetSurvey(ctx iris.Context) {
	survey_id := ctx.Params().Get("id")
	sess, ok := ctx.Values().Get("surveydb").(db.Session)
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

	var results []map[string]interface {}
	err := iter.All(&results);
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "failed to fetch survey",
		})
		return
	}
	if len(results) == 0 {
		ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
			"message": "no survey found",
		})
		return
	}

	survey := iris.Map{
		"id": results[0]["survey_id"],
		"title": results[0]["title"],
		"description": results[0]["description"],
		"questions": make([]iris.Map, 0),
	}

	if results[0]["question_id"] != nil {
		question := iris.Map{}
		for _, row := range results {
			question_id := row["question_id"]
			if question["id"] == nil {
				question = iris.Map{
					"id": question_id,
					"type": row["type"],
					"statement": row["statement"],
					"answers": make([]string, 0),
				}
			} else if question["id"] != question_id {
				survey["questions"] = append(
					survey["questions"].([]iris.Map),
					question,
				)
				question = iris.Map{
					"id": question_id,
					"type": row["type"],
					"statement": row["statement"],
					"answers": make([]string, 0),
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
