<html>

<head>
    <title>{{.Title}}</title>
</head>

<body>
    <form action="/category" method="post">
        <div>
            <input type="text" name="cateName" />
            <input type="submit" value="Add" />
        </div>
        <div>
            <table>
                <tr>
                    <th>Index</th>
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
</body>

</html>
