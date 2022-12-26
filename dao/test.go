package dao

func (d *Dao) CreateTestDemo() error {
	sql := "insert into user(id,name) values(2,'tom')"
	return d.mysql.Exec(sql).Error
}
