{{template "header" .}} {{template "nav" .}}
<form action="/user" method="post">
    <input type="text" hidden="hidden" name="pageIndex" value="{{.Pagination.Index}}" />
    <div class="wrap-primary">
        <label for="inputName" class="control-label">Name</label>
        <input type="text" name="username" id="username" />
        <label for="inputName" class="control-label">Password</label>
        <input type="password" name="password" id="password" />
        <label for="inputName" class="control-label">Nick</label>
        <input type="text" name="nick" id="nick" />
        <label for="inputName" class="control-label">Host</label>
        <input type="text" name="hostname" id="hostname" />

        <input type="submit" value="Add" class="btn btn-primary" />
    </div>
    <div class="wrap-table">
        <table class="table table-striped table-bordered table-hover table-condensed">
            <tr>
                <th>#</th>
                <th>Name</th>
                <th>Nick</th>
                <th>LastLoginTime</th>
                <th>LastLoginIP</th>
                <th>C.Time</th>
                <th>C.By</th>
                <th colspan="2">Action</th>
            </tr>

            {{if .Users}} {{range .Users}}
            <tr>
                <td><span name="userid">{{.ID}}</span></td>
                <td><span name="username">{{.Username}}</span></td>
                <td><span name="usernick">{{.Nick}}</span></td>
                <td><span name="userltime">{{.LastLoginTime}}</span></td>
                <td><span name="userlip">{{.LastLoginIP}}</span></td>
                <td><span name="userctime">{{.Operation.CreatedTime}}</span></td>
                <td><span name="usercby">{{.Operation.CreatedBy}}</span></td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal" class="btn btn-link edit">Edit</td>
                <td><a href="/user?id={{.ID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td>
            </tr>
            {{end}}
            {{else}}
                <tr>
                    <td colspan="7">
                        <blockquote>
                        <p>No data found.</p>
                        </blockquote>
                    </td>
                </tr>
            {{end}}

        </table>
    </div>
</form>
{{template "footer" .}}
