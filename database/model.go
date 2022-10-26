package database

type Model struct {
  Table string
}

func (M Model) New() *Orm {
      return new(Orm).Init("database", M.Table)
}

func (M Model) Insert(data map[string]string) int64 {
  return new(Orm).Init("database", M.Table).Insert(data)
}

func (M Model) Delete() int64 {
    return new(Orm).Init("database", M.Table).Delete()
}

func (M Model) Update(data map[string]string) int64 {
      return new(Orm).Init("database", M.Table).Delete()
}

func (M Model) Field(fields string) *Orm {
      return new(Orm).Init("database", M.Table).Field(fields)
}

func (M Model) Where(conds ...interface{}) *Orm {
      return new(Orm).Init("database", M.Table).Where(conds...)
}

func (M Model) Group(fields ...string) *Orm {
      return new(Orm).Init("database", M.Table).Group(fields...)
}

func (M Model) Having(field string, opr string, criteria int) *Orm {
      return new(Orm).Init("database", M.Table).Having(field, opr, criteria)
}

func (M Model) Order(field string, sort string) *Orm {
      return new(Orm).Init("database", M.Table).Order(field, sort)
}

func (M Model) Page(page int) *Orm {
      return new(Orm).Init("database", M.Table).Page(page)
}

func (M Model) Limit(limit int) *Orm {
      return new(Orm).Init("database", M.Table).Limit(limit)
}

func (M Model) Select() []map[string]string {
      return new(Orm).Init("database", M.Table).Select()
}

func (M Model) Find() map[string]string {
      return new(Orm).Init("database", M.Table).Find()
}

func (M Model) Value(field string) string {
      return new(Orm).Init("database", M.Table).Value(field)
}

func (M Model) Values(field string) []string {
      return new(Orm).Init("database", M.Table).Values(field)
}

func (M Model) Columns(fields ...string) map[string]string {
      return new(Orm).Init("database", M.Table).Columns(fields...)
}

func (M Model) Sum(field string) int {
      return new(Orm).Init("database", M.Table).Sum(field)
}

func (M Model) Count() int {
      return new(Orm).Init("database", M.Table).Count()
}

func (M Model) Exist(primary string) bool {
      return new(Orm).Init("database", M.Table).Exist(primary)
}
