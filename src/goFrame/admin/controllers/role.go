package controllers

import (
	"goFrame/models/system"
	"goFrame/enums"
	"strings"
	"strconv"
	"fmt"
	"time"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Index() {
	beego.Debug("RoleController-Index")
	//需要权限控制
	this.checkAuthor()

	//页面里按钮权限控制
	//this.Data["canEdit"] = this.checkActionRoleor("MenuController", "Edit")
	//this.Data["canDelete"] = this.checkActionRoleor("MenuController", "Delete")

	////获取角色列表
	//roles := system.RoleListGrid(this.page, this.pageSize)
	////rolelist["rows"] = this.Json_encode(data)
	//this.Data["role_rows"] = roles

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["datagrid"] = "admin/" + this.controllerName + "/datagrid.html"
	this.display()
}

//DataList 角色列表
func (this *RoleController) RoleList() {
	//获取角色列表
	data := system.RoleListGrid(this.page, this.pageSize)
	//定义返回的数据结构
	this.jsonResult(enums.JRCodeSucc, "", data)
}

/*
	DataGrid 后台用户角色管理页 表格获取数据
 */
func (this *RoleController) RoleDataGrid() {
	beego.Debug("RoleDataGrid")
	//直接反序化获取json格式的requestbody里的值
	var params system.RoleQueryParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	//data, total := system.RoleList(this.page, this.pageSize)
	data, total := system.RolePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RoleController) Save() {
	m := system.Role{}
	o := orm.NewOrm()
	var err error
	//获取form里的值
	if err = this.ParseForm(&m); err != nil {
		beego.Error(err.Error())
		this.jsonResult(enums.JRCodeFailed, "获取数据失败!", m.Id)
	}
	if m.Id == 0 {
		m.CreateTime = time.Now().Unix()
		if _, err := o.Insert(&m); err != nil {
			beego.Error(err.Error())
			this.jsonResult(enums.JRCodeFailed, "添加失败!", m.Id)
		}
	} else {
		m.UpdateTime = time.Now().Unix()
		if _, err := o.Update(&m); err != nil {
			beego.Error(err.Error())
			this.jsonResult(enums.JRCodeFailed, "编辑失败!", m.Id)
		}
	}
	this.jsonResult(enums.JRCodeSucc, "角色添加成功!", m.Id)
}

func (this *RoleController) RoleEdit() {
	beego.Debug("RoleEdit")
	// Edit 添加 编辑 页面
	//如果是Post请求，则由Save处理
	if this.Ctx.Request.Method == "POST" {
		this.Save()
	}
	Id, _ := this.GetInt(":id", 0)
	role := &system.Role{}
	var err error

	if Id > 0 {
		role, err = system.RoleGetById(Id)
		if err != nil {
			this.pageError("数据无效，请刷新后重试")
		}
	} else {
		//添加用户时默认状态为启用
		role.Status = enums.Enabled
	}
	this.Data["m"] = role
	this.setTpl("admin/role/edit.html", "layout/admin/layout_pullbox.html")
}

/*
	后台逻辑删除角色
 */
func (this *RoleController) RoleDelete() {
	if this.Ctx.Request.Method == "POST" {
		userIds := this.GetString("ids")
		ids := make([]int, 0, len(userIds))
		for _, id := range strings.Split(userIds, ",") {
			if id, err := strconv.Atoi(id); err == nil {
				ids = append(ids, id)
			} else {
				this.jsonResult(enums.JRCodeFailed, err.Error(), nil)
			}
		}
		query := orm.NewOrm().QueryTable(system.TableName("role"))
		//物理删除
		//if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		//逻辑删除
		if num, err := query.Filter("id__in", ids).Update(orm.Params{"logic_delete": 1}); err == nil {
			this.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 个角色!", num), nil)
		} else {
			beego.Error(err.Error())
			this.jsonResult(enums.JRCodeFailed, "角色删除失败!", nil)
		}
	} else {
		this.jsonResult(enums.JRCodeSucc, "请求方式错误!", nil)
	}
}
