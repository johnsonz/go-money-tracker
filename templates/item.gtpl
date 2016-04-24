{{template "header" .}} {{template "nav" .}}
<form action="/item" method="POST" enctype="multipart/form-data">
    <div class="wrap-field">
        <table class="table table-condensed">
            <tr>
                <td>Category:</td>
                <td>
                    <select name="category" id="category">
                        {{range .Categories}}
                        <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                        {{end}}
                    </select>
                </td>
            </tr>
            <tr>
                <td>Subcategory:</td>
                <td>
                    <select name="subcategory" id="subcategory">
                        {{range .Subcategories}}
                        <option value="{{.ID}}" {{if .Selected}}selected="selected" {{end}}>{{.Name}}</option>
                        {{end}}
                    </select>
                </td>
            </tr>
            <tr>
                <td>PurchasedDate:</td>
                <td>
                    <input type="text" name="purchaseddate" id="datepicker" />
                </td>
            </tr>
            <tr>
                <td>Store:</td>
                <td>
                    <input type="text" name="store" />
                </td>
            </tr>
            <tr>
                <td>Address:</td>
                <td>
                    <input type="text" name="address" />
                </td>
            </tr>
            <tr>
                <td>Remark:</td>
                <td>
                    <input type="text" name="remark" />
                </td>
            </tr>
            <tr>
                <td>Receipt Image:</td>
                <td>
                    <input type="file" name="receiptimage" />
                </td>
            </tr>
            <tr>
                <td colspan="2">
                    <input type="submit" value="Add" class="btn btn-primary" />
                </td>
            </tr>
        </table>

    </div>
    <div class="wrap-table">
        <table class="table table-striped table-bordered table-hover table-condensed">
            <tr>
                <th>Index</th>
                <th>Category</th>
                <th>Subcategory</th>
                <th>Store</th>
                <th>Address</th>
                <th>PurchasedDate</th>
                <th>Amount</th>
                <th>Receipt</th>
                <th>Remark</th>
                <th colspan="4">Action</th>
            </tr>
            {{range .Items}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Subcategory.Category.Name}}</td>
                <td>{{.Subcategory.Name}}</td>
                <td>{{.Store}}</td>
                <td>{{.Address}}</td>
                <td>{{.PurchasedDate}}</td>
                <td>{{.Amount}}</td>
                <td>{{if .Receipt}}<img class="smallimg" src="data:image/jpg;base64,{{.Receipt}}">{{else}}None{{end}}</td>
                <td>{{.Remark}}</td>
                <td><a href="/item?id={{.ID}}" class="btn btn-link">View</td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal" class="btn btn-link edit">Edit</td>
                <td><a href="/detail?id={{.ID}}" class="btn btn-link">Add</td>
                <td><a href="/item?id={{.ID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td>
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
                    <a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{.Pagination.Previous}}'" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{minus .Pagination.Index 1}}'">{{minus .Pagination.Index 2}}</a></li>
                {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                <!-- <li><a href="/item?page={{minus .Pagination.Index 1}}">{{minus .Pagination.Index 1}}</a></li> -->
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{minus .Pagination.Index 1}}'">{{minus .Pagination.Index 1}}</a></li>
                {{end}}
                <li class="active"><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{.Pagination.Index}}'">{{.Pagination.Index}}</a></li>
                {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{plus .Pagination.Index 1}}'">{{plus .Pagination.Index 1}}</a></li>
                {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                <li><a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{plus .Pagination.Index 2}}'">{{plus .Pagination.Index 2}}</a></li>
                {{end}} {{if le .Pagination.Next .Pagination.Count}}
                <li>
                    <a href="javascript:" onclick="this.href='/item?cid='+$('#category').val()+'&sid='+$('#subcategory').val()+'&page={{.Pagination.Next}}'" aria-label="Next">
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
