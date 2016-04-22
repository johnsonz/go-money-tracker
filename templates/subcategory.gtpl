{{template "header" .}} {{template "nav" .}}
<form action="subcategory" method="POST">
    <div class="wrap-primary">
        <label class="control-label">Category</label>
        <select name="category">
            {{range .Categories}}
            <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
            {{end}}
        </select>
        <label class="control-label">Subcategory</label>
        <input type="text" name="subcateName" />
        <input type="submit" value="Add" class="btn btn-primary" />
    </div>
    <div class="wrap-table">
        <table class="table table-striped table-bordered table-hover table-condensed">
            <tr>
                <th>#</th>
                <th>Subcategory</th>
                <th>CreatedTime</th>
                <th>CreatedBy</th>
            </tr>
            {{if .Subcategories}}{{range .Subcategories}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Name}}</td>
                <td>{{.CreatedTime}}</td>
                <td>{{.CreatedBy}}</td>
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
        <nav>
            <ul class="pagination">
                {{if le .Pagination.Previous 0}}
                <li class="disabled">
                    <a href="" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                {{else}}
                <li>
                    <a href="/subcategory?id={{.CategoryId}}&page={{.Pagination.Previous}}" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                <li><a href="/subcategory?id={{.CategoryId}}&page={{minus .Pagination.Index 2}}">{{minus .Pagination.Index 2}}</a></li>
                {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                <li><a href="/subcategory?id={{.CategoryId}}&page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li>
                {{end}}
                <li class="active"><a href="/subcategory?id={{.CategoryId}}&page={{.Pagination.Index}}">{{.Pagination.Index}}</a></li>
                {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                <li><a href="/subcategory?id={{.CategoryId}}&page={{plus .Pagination.Index 1}}">{{plus .Pagination.Index 1}}</a></li>
                {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                <li><a href="/subcategory?id={{.CategoryId}}&page={{plus .Pagination.Index 2}}">{{plus .Pagination.Index 2}}</a></li>
                {{end}} {{if le .Pagination.Next .Pagination.Count}}
                <li>
                    <a href="/subcategory?id={{.CategoryId}}&page={{.Pagination.Next}}" aria-label="Next">
                        <span aria-hidden="true">&raquo;</span>
                    </a>
                </li>
                {{else}}
                <li class="disabled">
                    <a href="" aria-label="Next">
                        <span aria-hidden="true">&raquo;</span>
                    </a>
                </li>
                {{end}}

            </ul>
        </nav>
    </div>

</form>
{{template "footer" .}}
