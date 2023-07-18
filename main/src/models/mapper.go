package model

import "xgo/main/src/gen"

func CompanyToProto(c *Company) *gen.Company {
	return &gen.Company{
		Id:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Employee:    c.EmployeeCount,
		Registered:  c.Registered,
		Type:        gen.Type(gen.Type_value[c.Type]),
	}
}

func CompanyFromProto(c *gen.Company) *Company {
	return &Company{
		ID:            c.Id,
		Name:          c.Name,
		Description:   c.Description,
		EmployeeCount: c.Employee,
		Registered:    c.Registered,
		Type:          gen.Type_name[int32(c.Type)],
	}
}
