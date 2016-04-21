{{template "header" .}} {{template "nav" .}}
<form action="/category" method="post">
    <div class="wrap-primary">
        <label for="inputCateName" class="control-label">Category</label>
        <input type="text" name="cateName" id="inputCateName" />
        <input type="submit" value="Add" class="btn btn-primary" />
    </div>
    <div class="wrap-table">
        <table class="table table-striped table-bordered table-hover table-condensed">
            <tr>
                <th>#</th>
                <th>Category</th>
                <th>Created Time</th>
                <th>Created By</th>
                <th colspan="2">Action</th>
            </tr>

            {{range .Categories}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Name}}</td>
                <td>{{.CreatedTime}}</td>
                <td>{{.CreatedBy}}</td>
                <td><a href="/category?id={{.ID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td>
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
                        <a href="/category?page={{.Pagination.Previous}}" aria-label="Previous">
                            <span aria-hidden="true">&laquo;</span>
                        </a>
                    </li>
                    {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                    <li><a href="/category?page={{minus .Pagination.Index 2}}">{{minus .Pagination.Index 2}}</a></li>
                    {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                    <li><a href="/category?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li>
                    {{end}}
                    <li class="active"><a href="/category?page={{.Pagination.Index}}">{{.Pagination.Index}}</a></li>
                    {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                    <li><a href="/category?page={{plus .Pagination.Index 1}}">{{plus .Pagination.Index 1}}</a></li>
                    {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                    <li><a href="/category?page={{plus .Pagination.Index 2}}">{{plus .Pagination.Index 2}}</a></li>
                    {{end}}  {{if le .Pagination.Next .Pagination.Count}}
                    <li>
                        <a href="/category?page={{.Pagination.Next}}" aria-label="Next">
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
</form>
{{template "footer" .}}
