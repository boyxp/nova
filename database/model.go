package database

import "strings"

type Model struct {
  Table string
}

func (M Model) New() *Orm {
      dbtag := "database"
      table := M.Table
      index := strings.Index(table, ".")
      if index>0 {
            dbtag = table[0:index]
            table = table[index+1:]
      }

      return new(Orm).Init(dbtag, table)
}

func (M Model) Insert(data map[string]string) int64 {
  return M.New().Insert(data)
}

func (M Model) Delete() int64 {
    return M.New().Delete()
}

func (M Model) Update(data map[string]string) int64 {
      return M.New().Update(data)
}

func (M Model) Field(fields string) *Orm {
      return M.New().Field(fields)
}

func (M Model) Where(conds ...interface{}) *Orm {
      return M.New().Where(conds...)
}

func (M Model) Group(fields ...string) *Orm {
      return M.New().Group(fields...)
}

func (M Model) Having(field string, opr string, criteria int) *Orm {
      return M.New().Having(field, opr, criteria)
}

func (M Model) Order(field string, sort string) *Orm {
      return M.New().Order(field, sort)
}

func (M Model) Page(page int) *Orm {
      return M.New().Page(page)
}

func (M Model) Limit(limit int) *Orm {
      return M.New().Limit(limit)
}

func (M Model) Select() []map[string]string {
      return M.New().Select()
}

func (M Model) Find(primary ...string) map[string]string {
      return M.New().Find(primary...)
}

func (M Model) Value(field string) string {
      return M.New().Value(field)
}

func (M Model) Values(field string) []string {
      return M.New().Values(field)
}

func (M Model) Columns(fields ...string) map[string]string {
      return M.New().Columns(fields...)
}

func (M Model) Map(fields ...string) map[string]map[string]string {
      return M.New().Map(fields...)
}

func (M Model) Sum(field string) int {
      return M.New().Sum(field)
}

func (M Model) Count() int {
      return M.New().Count()
}

func (M Model) Total() int {
      return M.New().Count()
}

func (M Model) TotalPage() int {
      return M.New().TotalPage()
}

func (M Model) Exist(primary string) bool {
      return M.New().Exist(primary)
}

func (M Model) Relate(list *[]map[string]string, fields string) {
      M.New().Relate(list, fields)
}

