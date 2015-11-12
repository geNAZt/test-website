{{ template "layout.tpl" . }}

{{ define "content" }}
<div id="inner-content" class="col-xs-4 col-sm-4">
    <div class="box">
        <div class="box-content">
            <form action="/profile/settings/" method="POST" enctype="multipart/form-data">
                <div class="text-center">
                    <h3 class="page-header">MineTracker Settings</h3>
                </div>
                <div class="form-group{{ if isset .errors "avatar" }} has-error{{ end }}">
                    <label class="control-label">Avatar</label>
                    <br>
                    <div style="display: inline-table;">
                        <img src="{{ .user.Avatar | getAvatar }}" class="img-circle" alt="avatar">
                    </div>
                    <div style="display: inline-table; position: relative; top: -10px; margin-right: 128px;">
                        Upload new one (40x40 and max 16kb) <input type="file" name="avatarfile">
                    </div>
                    <div style="display: inline-table;">
                        or use Gravatar <input type="checkbox" name="avatargravatar" value="true">
                    </div>
                    {{ if isset .errors "avatar" }}<p class="txt-danger">{{ .errors.avatar }}</p>{{ end }}
                </div>
                <div class="form-group{{ if isset .errors "display" }} has-error{{ end }}">
                    <label class="control-label">Displayname</label>
                    <input type="text" class="form-control" name="display">
                    {{ if isset .errors "display" }}<p class="txt-danger">{{ .errors.display }}</p>{{ end }}
                </div>
                <div class="form-group{{ if isset .errors "password" }} has-error{{ end }}">
                    <label class="control-label">New password</label>
                    <input type="password" class="form-control" name="pass1">
                    <input type="password" class="form-control" name="pass2">
                    {{ if isset .errors "password" }}<p class="txt-danger">{{ .errors.password }}</p>{{ end }}
                </div>
                <div class="text-center">
                    <button type="submit" class="btn btn-primary">Save changes</button>
                </div>
            </form>
        </div>
    </div>
</div>
{{ end }}

{{ define "js" }}

{{ end }}

{{ define "css" }}

{{ end }}