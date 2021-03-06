{{template "header" .}} {{template "nav" .}}
<form action="/user" method="post">
    <input type="text" hidden="hidden" name="pageIndex" value="{{.Pagination.Index}}" />
    <div class="wrap-primary">
        <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#Modal-Create-User">Create new user</button>
    </div>
    <div class="wrap-table">
        <table class="table table-striped table-bordered table-hover table-condensed">
            <tr>
                <th>#</th>
                <th>Name</th>
                <th>Nick</th>
                <th>Host</th>
                <th>LastLoginTime</th>
                <th>LastLoginIP</th>
                <th>C.Time</th>
                <th>C.By</th>
                <th colspan="2">Action</th>
            </tr>

            {{if .Users}} {{range $index,$user:=.Users}}
            <tr>
                <td hidden="hidden"><span name="userid">{{.ID}}</span></td>
                <td><span name="userindex">{{plus $index 1}}</span></td>
                <td><span name="username">{{.Username}}</span></td>
                <td><span name="usernick">{{.Nick}}</span></td>
                <td><span name="userhost">{{.Hostname}}</span></td>
                <td><span name="userltime">{{.LastLoginTime}}</span></td>
                <td><span name="userlip">{{.LastLoginIP}}</span></td>
                <td><span name="userctime">{{.Operation.CreatedTime}}</span></td>
                <td><span name="usercby">{{.Operation.CreatedBy}}</span></td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal-User" class="btn btn-link edit">Edit</td>
                <!-- <td><a href="/user?id={{.ID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td> -->
                <td><a href="javascript:" class="btn btn-link" data-toggle="modal" data-target="#Modal-Confirm-User">Delete</td>

            </tr>
            {{end}}
            {{else}}
                <tr>
                    <td colspan="9">
                        <blockquote>
                        <p>No data found.</p>
                        </blockquote>
                    </td>
                </tr>
            {{end}}

        </table>
        <nav>
            <ul class="pagination">
                {{if le .Pagination.Previous 0}}
                <li class="disabled">
                    <a href="#" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                    </li>
                    {{else}}
                    <li>
                        <a href="/user?page={{.Pagination.Previous}}" aria-label="Previous">
                            <span aria-hidden="true">&laquo;</span>
                        </a>
                    </li>
                    {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                    <li><a href="/user?page={{minus .Pagination.Index 2}}">{{minus .Pagination.Index 2}}</a></li>
                    {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                    <li><a href="/user?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li>
                    {{end}}
                    <li class="active"><a href="/user?page={{.Pagination.Index}}" id="pageIndex">{{.Pagination.Index}}</a></li>
                    {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                    <li><a href="/user?page={{plus .Pagination.Index 1}}">{{plus .Pagination.Index 1}}</a></li>
                    {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                    <li><a href="/user?page={{plus .Pagination.Index 2}}">{{plus .Pagination.Index 2}}</a></li>
                    {{end}} {{if le .Pagination.Next .Pagination.Count}}
                    <li>
                        <a href="/user?page={{.Pagination.Next}}" aria-label="Next">
                            <span aria-hidden="true">&raquo;</span>
                        </a>
                    </li>
                    {{else}}
                    <li class="disabled">
                        <a href="#" aria-label="Next">
                            <span aria-hidden="true">&raquo;</span>
                        </a>
                    </li>
                    {{end}}

                    </ul>
                    </nav>
    </div>
    <div class="modal fade" id="Modal-Confirm-User" tabindex="-1" role="dialog" aria-labelledby="ModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    Confirm
                </div>
                <div class="modal-body">
                    This record will be permanently deleted and cannot be recovered. Are you sure?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
                    <a class="btn btn-danger btn-ok userdel">Delete</a>
                </div>
            </div>
        </div>
    </div>
    <input type="text" name="updatedid" id="updatedid" hidden="hidden">
    <div class="modal fade" id="Modal-User" tabindex="-1" role="dialog" aria-labelledby="ModalLabel">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="ModalLabel">Edit</h4>
                </div>
                <div class="modal-body">
                        <div class="form-group">
                            <label for="username" class="control-label">Username:</label>
                            <input type="text" name="updatedname" class="form-control" id="username">
                        </div>
                        <div class="form-group">
                            <label for="userpassword" class="control-label">Password:</label>
                            <input type="text" name="updatedpassword" class="form-control" id="userpassword">
                        </div>
                        <div class="form-group">
                            <label for="usernick" class="control-label">Nick:</label>
                            <input type="text" name="updatednick" class="form-control" id="usernick">
                        </div>
                        <div class="form-group">
                            <label for="userhost" class="control-label">Hostname:</label>
                            <input type="text" name="updatedhost" class="form-control" id="userhost">
                        </div>
                        <div class="form-group">
                            <label for="createdtime" class="control-label">Created Time:</label>
                            <input type="text" class="form-control" id="createdtime" disabled="disabled">
                            </textarea>
                        </div>
                        <div class="form-group">
                            <label for="createdby" class="control-label">Created By:</label>
                            <input type="text" class="form-control" id="createdby" disabled="disabled">
                            </textarea>
                        </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
                    <input type="submit" name="update" value="Update" class="btn btn-primary" />
                </div>
            </div>
        </div>
    </div>
    <div class="modal fade" id="Modal-Create-User" tabindex="-1" role="dialog" aria-labelledby="ModalLabel">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="ModalLabel">Create new user</h4>
                </div>
                <div class="modal-body">
                        <div class="form-group">
                            <label for="createdname" class="control-label">Username:</label>
                            <input type="text" name="createdname" class="form-control" id="createdname">
                        </div>
                        <div class="form-group">
                            <label for="createdpassword" class="control-label">Password:</label>
                            <input type="text" name="createdpassword" class="form-control" id="createdpassword">
                        </div>
                        <div class="form-group">
                            <label for="creatednick" class="control-label">Nick:</label>
                            <input type="text" name="creatednick" class="form-control" id="creatednick">
                        </div>
                        <div class="form-group">
                            <label for="createdhost" class="control-label">Hostname:</label>
                            <input type="text" name="createdhost" class="form-control" id="createdhost">
                        </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
                    <input type="submit" name="create" value="Create" class="btn btn-primary" />
                </div>
            </div>
        </div>
    </div>
</form>
{{template "footer" .}}
