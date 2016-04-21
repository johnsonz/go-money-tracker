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
                <th colspan="3">Action</th>
            </tr>

            {{range .Categories}}
            <tr>
                <td><span name="cateid">{{.ID}}</span></td>
                <td><span name="catename">{{.Name}}</span></td>
                <td><span name="catectime">{{.CreatedTime}}</span></td>
                <td><span name="catecby">{{.CreatedBy}}</span></td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal" class="btn btn-link edit">Edit</td>
                <td><a href="/subcategory?id={{.ID}}" class="btn btn-link">Add</td>
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
                    {{end}} {{if le .Pagination.Next .Pagination.Count}}
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
    <div class="modal fade" id="Modal" tabindex="-1" role="dialog" aria-labelledby="ModalLabel">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title" id="ModalLabel">Edit</h4>
      </div>
      <div class="modal-body">
        <form>
          <div class="form-group">
            <label for="catename" class="control-label">Name:</label>
            <input type="text" class="form-control" id="catename">
          </div>
          <div class="form-group">
            <label for="createdtime" class="control-label">Created Time:</label>
            <input type="text" class="form-control" id="createdtime" disabled="disabled"></textarea>
          </div>
          <div class="form-group">
            <label for="createdby" class="control-label">Created By:</label>
            <input type="text" class="form-control" id="createdby" disabled="disabled"></textarea>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
        <button type="button" class="btn btn-primary">Update</button>
      </div>
    </div>
  </div>
</div>
</form>
{{template "footer" .}}
