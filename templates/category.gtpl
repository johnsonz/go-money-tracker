{{template "header" .}}
{{template "nav" .}}
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
            </tr>
            {{range .Categories}}
            <tr>

                <td>{{.ID}}</td>
                <td>{{.Name}}</td>
                <td>{{.CreatedTime}}</td>
                <td>{{.CreatedBy}}</td>

            </tr>
            {{end}}
        </table>
        <nav>
  <ul class="pagination">
    <li>
      <a href="#" aria-label="Previous">
        <span aria-hidden="true">&laquo;</span>
      </a>
    </li>
    <li><a href="#">1</a></li>
    <li><a href="#">2</a></li>
    <li><a href="#">3</a></li>
    <li><a href="#">4</a></li>
    <li><a href="#">5</a></li>
    <li>
      <a href="#" aria-label="Next">
        <span aria-hidden="true">&raquo;</span>
      </a>
    </li>
  </ul>
</nav>
        
    </div>
</form>
{{template "footer" .}}
