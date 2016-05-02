{{template "header" .}} {{template "nav" .}}
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
                <th>#</th>
                <th>Name</th>
                <th>Price</th>
                <th>Quantity</th>
                <th>Amount</th>
                <th>LabelOne</th>
                <th>LabelTwo</th>
                <th>Remark</th>
                <th>Created Time</th>
                <th>Created By</th>
                <th colspan="2">Action</th>
            </tr>
            {{if .Details}}{{range $index,$detail:=.Details}}
            <tr>
                <td><span name="detailid">{{plus $index 1}}</span></td>
                <td><span name="detailname">{{.Name}}</span></td>
                <td><span name="detailprice">{{.Price}}</span></td>
                <td><span name="detailquan">{{.Quantity}}</span></td>
                <td><span name="detailamount">{{getamount .Price .Quantity}}</span></td>
                <td><span name="detaillone">{{if .LabelOne}}<img class="smallimg" src="data:image/jpg;base64,{{.LabelOne}}">{{else}}None{{end}}</span></td>
                <td><span name="detailltwo">{{if .LabelTwo}}<img class="smallimg" src="data:image/jpg;base64,{{.LabelTwo}}">{{else}}None{{end}}</span></td>
                <td><span name="detailremark">{{.Remark}}</td>
                <td><span name="detailctime">{{.Operation.CreatedTime}}</span></td>
                <td><span name="detailcby">{{.Operation.CreatedBy}}</span></td>
                <td><a href="javascript:" data-toggle="modal" data-target="#Modal-Detail" class="btn btn-link edit">Edit</td>
                <!-- <td><a href="/detail?id={{.ID}}&iid={{$.ItemID}}&action=del&page={{$.Pagination.Index}}" class="btn btn-link">Delete</td> -->
                <td><a href="javascript:" class="btn btn-link" data-toggle="modal" data-target="#Modal-Confirm-Detail">Delete</td>
                </tr>
                {{end}}
                {{else}}
                    <tr>
                        <td colspan="11">
                            <blockquote>
                            <p>No data found.</p>
                            </blockquote>
                        </td>
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
                    <li class="active"><a href="javascript:" onclick="this.href='/detail?id={{.ItemID}}&page={{.Pagination.Index}}'" id="pageIndex">{{.Pagination.Index}}</a></li>
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
    <div class="modal fade" id="Modal-Confirm-Detail" tabindex="-1" role="dialog" aria-labelledby="ModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    Confirm
                </div>
                <div class="modal-body">
                    This record will be permanently deleted and cannot be recovered. Are you sure?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
                    <a class="btn btn-danger btn-ok detaildel">Delete</a>
                </div>
            </div>
        </div>
    </div>
    <input type="text" name="updatedid" id="updatedid" hidden="hidden">
    <div class="modal fade" id="Modal-Detail" tabindex="-1" role="dialog" aria-labelledby="ModalLabel">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="ModalLabel">Edit</h4>
                </div>
                <div class="modal-body">
                    <form>
                        <div class="form-group">
                            <label for="updatedname" class="control-label">Name:</label>
                            <input type="text" class="form-control" name="updatedname" id="updatedname">
                        </div>
                        <div class="form-group">
                            <label for="updatedprice" class="control-label">Price:</label>
                            <input type="number" step="0.01" min="0" class="form-control" name="updatedprice" id="updatedprice">
                        </div>
                        <div class="form-group">
                            <label for="updatedquantity" class="control-label">Quantity:</label>
                            <input type="number" step="1" min="0" name="updatedquantity" id="updatedquantity" class="datepicker form-control" />
                        </div>
                        <div class="form-group">
                            <label for="updatedamount" class="control-label">Amount:</label>
                            <input type="text" class="form-control" name="updatedamount" id="updatedamount" disabled="disable">
                        </div>
                        <div class="form-group">
                            <label for="updatedlone" class="control-label">LabelOne:</label>
                            <div id="wraplone"></div>
                            <input type="file" class="form-control" name="updatedlone" id="updatedlone">
                        </div>
                        <div class="form-group">
                            <label for="updatedltwo" class="control-label">LabelTwo:</label>
                            <div id="wrapltwo"></div>
                            <input type="file" class="form-control" name="updatedltwo" id="updatedltwo">
                        </div>
                        <div class="form-group">
                            <label for="updatedremark" class="control-label">Remark:</label>
                            <input type="text" class="form-control" name="updatedremark" id="updatedremark">
                        </div>
                        <div class="form-group">
                            <label for="createdtime" class="control-label">Created Time:</label>
                            <input type="text" class="form-control" id="createdtime" disabled="disabled">
                        </div>
                        <div class="form-group">
                            <label for="createdby" class="control-label">Created By:</label>
                            <input type="text" class="form-control" id="createdby" disabled="disabled">
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" name="cancel" data-dismiss="modal">Cancel</button>
                    <input type="submit" name="update" value="Update" class="btn btn-primary" />
                </div>
            </div>
        </div>
    </div>
</form>
{{template "footer" .}}
