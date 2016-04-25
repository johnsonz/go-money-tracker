{{template "header" .}}
{{template "nav" .}}
    <form action="/detail" method="POST" enctype="multipart/form-data">
        <input type="hidden" name="itemid" value="{{.ItemID}}" readonly="readonly" />
        <div class="wrap-field">
            <table class="table table-condensed">
                <tr>
                    <td>Name:</td>
                    <td>
                        <input type="text" name="name" />
                    </td>
                </tr>
                <tr>
                    <td>Price:</td>
                    <td>
                        <input type="number" step="0.01" min="0" name="price" />
                    </td>
                </tr>
                <tr>
                    <td>Quantity:</td>
                    <td>
                        <input type="number" min="0" name="quantity" value="1" />
                    </td>
                </tr>
                <tr>
                    <td>Label One:</td>
                    <td>
                        <input type="file" name="labelone" />
                    </td>
                </tr>
                <tr>
                    <td>Label Two:</td>
                    <td>
                        <input type="file" name="labeltwo" />
                    </td>
                </tr>
                <tr>
                    <td>Remark:</td>
                    <td>
                        <input type="text" name="remark" />
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
                    <th>Name</th>
                    <th>Price</th>
                    <th>Quantity</th>
                    <th>Amount</th>
                    <th>LabelOne</th>
                    <th>LabelTwo</th>
                    <th>Remark</th>
                    <th colspan="2">Action</th>
                </tr>
                {{range .Details}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.Price}}</td>
                    <td>{{.Quantity}}</td>
                    <td>{{getamount .Price .Quantity}}</td>
                    <td>{{if .LabelOne}}<img class="smallimg" src="data:image/jpg;base64,{{.LabelOne}}">{{else}}None{{end}}</td>
                    <td>{{if .LabelTwo}}<img class="smallimg" src="data:image/jpg;base64,{{.LabelTwo}}">{{else}}None{{end}}</td>
                    <td>{{.Remark}}</td>
                    <td><a href="javascript:" data-toggle="modal" data-target="#Modal-Detail" class="btn btn-link edit">Edit</td>
                    <td><a href="/detail?id={{.ID}}&iid={{$.ItemID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td>
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
                        <a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{.Pagination.Previous}}'" aria-label="Previous">
                            <span aria-hidden="true">&laquo;</span>
                        </a>
                    </li>
                    {{end}} {{if gt (minus .Pagination.Index 2) 0}}
                    <li><a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{minus .Pagination.Index 1}}'">{{minus .Pagination.Index 2}}</a></li>
                    {{end}} {{if gt (minus .Pagination.Index 1) 0}}
                    <li><a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{minus .Pagination.Index 1}}'">{{minus .Pagination.Index 1}}</a></li>
                    {{end}}
                    <li class="active"><a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{.Pagination.Index}}'">{{.Pagination.Index}}</a></li>
                    {{if le (plus .Pagination.Index 1) .Pagination.Count}}
                    <li><a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{plus .Pagination.Index 1}}'">{{plus .Pagination.Index 1}}</a></li>
                    {{end}} {{if le (plus .Pagination.Index 2) .Pagination.Count}}
                    <li><a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{plus .Pagination.Index 2}}'">{{plus .Pagination.Index 2}}</a></li>
                    {{end}} {{if le .Pagination.Next .Pagination.Count}}
                    <li>
                        <a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{.Pagination.Next}}'" aria-label="Next">
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
