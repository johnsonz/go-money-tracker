<html>
<head>
    <title></title>
</head>
<body>
    <form action="subcategory" method="POST">
<div>
    Category:<select name="category">
        {{range .Categories}}
        <option value="{{.ID}}" {{if .Selected}}selected="selected"{{end}}>{{.Name}}</option>
        {{end}}
    </select>
    Subcategory:<input type="text" name="subcateName"/>
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
