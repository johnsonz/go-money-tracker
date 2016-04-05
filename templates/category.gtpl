<html>

<head>
    <title>{{.Title}}</title>
</head>

<body>
    <form action="/category" method="post">
        <div>
            <input type="text" />
            <input type="submit" value="Add" />
        </div>
        <div>
            <table>
                <tr>
                    <th>
                        Index
                    </th>
                    <th>
                        Category
                    </th>
                    <th>
                        Created Time
                    </th>
                    <th>
                        Created By
                    </th>
                </tr>
                <tr>
                    {{range .Categorys}}
                    <td>
                        {{.ID}}
                    </td>
                    <td>
                        {{.Name}}
                    </td>
                    <td>
                        {{.CreatedTime}}
                    </td>
                    <td>
                        {{.Createdby}}
                    </td>
                    {{end}}
                </tr>
            </table>
        </div>
    </form>
</body>

</html>
