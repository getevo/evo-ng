package intl

import (
	"fmt"
	"github.com/getevo/evo-ng"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

var dbo *gorm.DB

type Model struct {
	LanguageCode language.Tag `gorm:"-"`
}

type ModelTranslation struct {
	Table       string `gorm:"primaryKey"`
	Field       string `gorm:"primaryKey"`
	Key         string `gorm:"primaryKey"`
	Language    string `gorm:"primaryKey"`
	Translation string
}

func (ModelTranslation) TableName() string {
	return "model_translation"
}

var modelLocalization = false

func EnableModelLocalization() {
	if modelLocalization {
		return
	}
	dbo = evo.GetDBO()
	modelLocalization = true
	var controller Model
	dbo.AutoMigrate(&ModelTranslation{})

	dbo.Callback().Create().After("*").Register("l10n:create", controller.onCreateOrUpdate)
	dbo.Callback().Delete().After("*").Register("l10n:delete", controller.onDelete)
	dbo.Callback().Update().After("*").Register("l10n:update", controller.onCreateOrUpdate)
	dbo.Callback().Query().After("gorm:query").Register("l10n:select", controller.onSelect)
	dbo.Callback().Row().After("gorm:query").Register("l10n:select", controller.onSelect)

}

func (n Model) onCreateOrUpdate(db *gorm.DB) {
	if _, skip := db.Statement.Get("l10n:skip"); !skip {
		if db.Error == nil && db.Statement.Schema != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Struct:
				if db.Statement.Schema.Name == "ModelTranslation" {
					return
				}
				var pk = ""
				var lang language.Tag

				if l, ok := db.Statement.Get("lang"); ok {
					lang = GuessLocale(fmt.Sprint(l))
				} else {
					if field := db.Statement.Schema.LookUpField("LanguageCode"); field != nil {
						if fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue); !isZero {
							if v, ok := fieldValue.(language.Tag); ok {
								lang = v
							}
						}
					} else {
						lang = GuessLocale("")
					}
					db.Statement.Set("lang", lang)
				}

				for _, field := range db.Statement.Schema.PrimaryFields {
					if fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue); !isZero {
						pk += "#" + fmt.Sprint(fieldValue)
					}
				}

				for _, field := range db.Statement.Schema.Fields {
					if field.Name == "LanguageCode" {
						field.Set(db.Statement.ReflectValue, lang)
					}
					if _, ok := field.Tag.Lookup("l10n"); ok {
						if fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue); !isZero {
							var translation = ModelTranslation{
								Table:       db.Statement.Schema.Table,
								Field:       field.DBName,
								Key:         pk,
								Language:    lang.String(),
								Translation: fmt.Sprint(fieldValue),
							}

							dbo.Clauses(clause.OnConflict{
								DoUpdates: clause.AssignmentColumns([]string{"translation"}), // column needed to be updated
							}).Create(&translation)

							for _, t := range supported {
								if t.String() != lang.String() {
									translation.Language = t.String()
									dbo.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&translation)
								}

							}

						}
					}
				}

			}

		}
	}
}

func (n Model) onDelete(db *gorm.DB) {
	if _, skip := db.Statement.Get("l10n:skip"); !skip {
		if db.Error == nil && db.Statement.Schema != nil {
			var pk = ""
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Struct:
				for _, field := range db.Statement.Schema.PrimaryFields {
					if fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue); !isZero {
						pk += "#" + fmt.Sprint(fieldValue)
					}
				}

				for _, field := range db.Statement.Schema.Fields {
					if _, ok := field.Tag.Lookup("l10n"); ok {

						var translation = ModelTranslation{
							Table: db.Statement.Schema.Table,
							Key:   pk,
						}
						dbo.Delete(&translation)

					}
				}

			}

		}
	}
}

func (n Model) onSelect(db *gorm.DB) {
	if _, skip := db.Statement.Get("l10n:skip"); !skip {
		if db.Error == nil && db.Statement.Schema != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice:
				var lang language.Tag
				if l, ok := db.Statement.Get("lang"); ok {
					lang = GuessLocale(fmt.Sprint(l))
				} else {
					lang = GuessLocale("")
					db.Statement.Set("lang", lang)
				}
				var langValue = reflect.ValueOf(Model{lang})
				for idx := 0; idx < db.Statement.ReflectValue.Len(); idx++ {
					obj := db.Statement.ReflectValue.Index(idx)
					var pk = ""
					for i, _ := range db.Statement.Schema.PrimaryFields {
						pk += "#" + fmt.Sprint(obj.Field(i).Interface())
					}
					for i, field := range db.Statement.Schema.Fields {
						if field.Name == "LanguageCode" {
							obj.Field(i).Set(langValue)
							continue
						}
						if _, ok := field.Tag.Lookup("l10n"); ok && field.DataType == schema.String {

							var translation ModelTranslation
							if dbo.Where(ModelTranslation{
								Table:    db.Statement.Schema.Table,
								Key:      pk,
								Field:    field.DBName,
								Language: lang.String(),
							}).Take(&translation).RowsAffected != 0 {
								obj.Field(i).SetString(translation.Translation)
							} else {
								//fix missing translation
								var translation = ModelTranslation{
									Table:       db.Statement.Schema.Table,
									Field:       field.DBName,
									Key:         pk,
									Language:    lang.String(),
									Translation: "",
								}
								dbo.Clauses(clause.OnConflict{
									DoUpdates: clause.AssignmentColumns([]string{"translation"}), // column needed to be updated
								}).Create(&translation)
							}

							obj.Field(i)
						}
					}
				}

			case reflect.Struct:
				var pk = ""
				var lang language.Tag
				if l, ok := db.Statement.Get("lang"); ok {
					lang = GuessLocale(fmt.Sprint(l))
				} else {
					lang = GuessLocale("")
					db.Statement.Set("lang", lang)
				}

				for _, field := range db.Statement.Schema.PrimaryFields {
					if fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue); !isZero {
						pk += "#" + fmt.Sprint(fieldValue)
					}
				}

				for _, field := range db.Statement.Schema.Fields {
					if field.Name == "LanguageCode" {
						field.Set(db.Statement.ReflectValue, lang)
						continue
					}
					if _, ok := field.Tag.Lookup("l10n"); ok && field.DataType == schema.String {
						var translation ModelTranslation
						if dbo.Where(ModelTranslation{
							Table:    db.Statement.Schema.Table,
							Key:      pk,
							Field:    field.DBName,
							Language: lang.String(),
						}).Take(&translation).RowsAffected != 0 {
							field.Set(db.Statement.ReflectValue, translation.Translation)
						} else {
							//fix missing translation
							var translation = ModelTranslation{
								Table:       db.Statement.Schema.Table,
								Field:       field.DBName,
								Key:         pk,
								Language:    lang.String(),
								Translation: "",
							}
							dbo.Clauses(clause.OnConflict{
								DoUpdates: clause.AssignmentColumns([]string{"translation"}), // column needed to be updated
							}).Create(&translation)
						}

					}
				}

			}

		}
	}
}
