<html>

<head>
    <title>Subcategory</title>
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap.min.css">
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="static/css/style.css">
</head>

<body>
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
                {{range .Subcategories}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.CreatedTime}}</td>
                    <td>{{.CreatedBy}}</td>
                </tr>
                {{end}}
            </table>
        </div>

    </form>
    <script src="/static/js/jquery-1.12.3.min.js"></script>
    <script src="static/bootstrap-3.3.6/js/bootstrap.min.js"></script>
</body>

</html>
