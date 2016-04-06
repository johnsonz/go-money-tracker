<html>
<head>
    <title></title>
</head>
<body>
    <form action="subcategory" method="POST">
<div>
    <select name="category">
        {{range .Categories}}
        <option value="{{.ID}}">{{.Name}}</option>
        {{end}}
    </select>
    <input type="text" name="subcateName"/>
    <input type="submit" value="Add" />
</div>
        <div>
            <table>
                <tr>
                    <th>Index</th>
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
</body>
</html>
