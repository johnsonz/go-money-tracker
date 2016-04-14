<html>

<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap.min.css">
    <link rel="stylesheet" href="static/bootstrap-3.3.6/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="static/css/style.css">
</head>

<body>
    <form action="/category" method="post">
        <div class="wrap-primary">
            <label for="inputCateName" class="control-label">Category</label>
            <input type="text" name="cateName" id="inputCateName"/>
            <input type="submit" value="Add" class="btn btn-primary"/>
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
        </div>
    </form>
    <script src="/static/js/jquery-1.12.3.min.js"></script>
    <script src="static/bootstrap-3.3.6/js/bootstrap.min.js"></script>
</body>

</html>
