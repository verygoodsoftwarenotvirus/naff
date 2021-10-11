package frontend

func favicon() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg"><style>path { fill: #666; }; @media (prefers-color-scheme: dark) { path { fill: #FFFFFF; } }</style>
    <circle cx="256" cy="256" r="256"/>
</svg>`
}

func loginAuthPartial() string {
	return `<div class="container">
    <div class="row">
        <div class="col-3"></div>
        <div class="col-6">
            <h1 class="h3 mb-3 text-center fw-normal">Log in</h1>
            <form hx-post="/auth/submit_login" hx-target="#content" hx-ext="json-enc, ajax-header, event-header">
                <div class="form-floating"><input id="usernameInput" required type="text" placeholder="username" minlength=4 name="username" placeholder="username" class="form-control"><label for="usernameInput">username</label></div>
                <div class="form-floating"><input id="passwordInput" required type="password" minlength=8 name="password" placeholder="password" class="form-control"><label for="passwordInput">password</label></div>
                <div class="form-floating"><input id="totpTokenInput" required type="text" pattern="\d{6}" minlength=6 maxlength=6 name="totpToken" placeholder="123456" class="form-control"><label for="totpTokenInput">2FA Token</label></div>
                <input type="hidden" name="redirectTo" value="{{ .RedirectTo }}" />
                <hr />
                <button id="loginButton" class="w-100 btn btn-lg btn-primary" type="submit">Log in</button>
            </form>
            <p class="text-center"><sub><a hx-target="#content" hx-push-url="/register" hx-get="/components/registration_prompt">Register instead</a></sub></p>
        </div>
        <div class="col-3"></div>
    </div>
</div>`
}

func registrationAuthPartial() string {
	return `<div class="container">
    <div class="row">
        <div class="col-3"></div>
        <div class="col-6">
            <h1 class="h3 mb-3 text-center fw-normal">Register</h1>
            <form hx-post="/auth/submit_registration" hx-ext="json-enc, ajax-header, event-header">
                <div class="form-floating"><input id="usernameInput" required type="text" placeholder="username" minlength=4 name="username" placeholder="username" class="form-control"><label for="usernameInput">username</label></div>
                <div class="form-floating"><input id="passwordInput" required type="password" minlength=8 name="password" placeholder="password" class="form-control"><label for="passwordInput">password</label></div>
                <hr />
                <button id="registrationButton" class="w-100 btn btn-lg btn-primary" type="submit">Register</button>
            </form>
            <p class="text-center"><sub><a hx-target="#content" hx-push-url="/login" hx-get="/components/login_prompt">Login instead</a></sub></p>
        </div>
        <div class="col-3"></div>
    </div>
</div>`
}

func registrationSuccessAuthPartial() string {
	return `<div class="container">
    <div class="row">
        <div class="col-3"></div>
        <div class="col-6">
            <h1 class="h3 mb-3 text-center fw-normal">Verify 2FA</h1>
            <form hx-post="/auth/verify_two_factor_secret" hx-ext="json-enc,ajax-header,event-header">
                <img id="twoFactorSecretQRCode" src={{ .TwoFactorQRCode }}>
                <div class="form-floating"><input id="totpTokenInput" required type="text" pattern="\d{6}" minlength=6 maxlength=6 name="totpToken" placeholder="123456" class="form-control"><label for="totpTokenInput">2FA Token</label></div>
                <input id="userID" type="hidden" name="userID" value="{{ .UserID }}" />
                <hr />
                <button id="totpTokenSubmitButton" class="w-100 btn btn-lg btn-primary" type="submit">Verify</button>
            </form>
        </div>
        <div class="col-3"></div>
    </div>
</div>`
}

func accountSettingsPartial() string {
	return `<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Account Settings</h1>
</div>
<div class="col-md-8 order-md-1">
    <div class="mb3">
        <form class="needs-validation" novalidate="">
        <label for="Name">Name</label>
        <div class="input-group">
            <input class="form-control" type="text" id="Name" placeholder="Name"required="" value="{{ .Account.Name }}" />
            <div class="invalid-feedback" style="width: 100%;">Name is required.</div>
        </div>

        <button class="btn btn-primary btn-lg btn-block mt-3" type="submit">Save</button>
        </form>
    </div>

    <hr class="mb-4" />

    <h3>Billing</h3>
    <div id="billing" class="mb3">
        <div>
            <select class="form-select form-select-lg mb-3" name="desiredPlan">
                <option selected>Please select a plan</option>
            </select>
            <button class="btn btn-primary btn-lg btn-block" type="submit" id="beginCheckoutButton">Checkout</button>
        </div>
    </div>

    <script type="text/javascript">
        // Create an instance of the Stripe object with your publishable API key
        let stripe = Stripe("pk_test_51IrCfgJ45Mr1esdKdRKAuAAH6U17SJFTeiCSWKqQzN5t8O3rbRBD5o1XjY2h5HG0hh0v4f3NHsaHC6KCp2NJPNm500MOMCpc7f");
        let checkoutButton = document.getElementById("beginCheckoutButton");

        checkoutButton.addEventListener("click", function () {
            fetch("/billing/checkout/begin", {
                method: "POST",
            })
            .then(function (response) {
                return response.json();
            })
            .then(function (response) {
                return stripe.redirectToCheckout({ sessionId: response.sessionID });
            })
            .then(function (result) {
                // If redirectToCheckout fails due to a browser or network
                // error, you should display the localized error message to your
                // customer using error.message.
                if (result.error) {
                    alert(result.error.message);
                }
            })
            .catch(function (error) {
                console.error("Error:", error);
            });
        });
    </script>
</div>`
}

func adminSettingsPartial() string {
	return `<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Admin Settings</h1>
</div>
<div class="col-md-8 order-md-1">
    <form class="needs-validation" novalidate="">
        <button class="btn btn-danger btn-lg btn-block" hx-confirm="This will log all users out, are you sure about that?" type="submit">Cycle Cookie Secret</button>
    </form>
</div>`
}

func userSettingsPartial() string {
	return `<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">User Settings</h1>
</div>
<div class="col-md-8 order-md-1">
    <form class="needs-validation" novalidate="">
        <div class="mb3">
            <label for="Username">Username</label>
            <div class="input-group">
                <input class="form-control" type="text" id="Username" placeholder="Username"required="" value="{{ .Username }}" />
                <div class="invalid-feedback" style="width: 100%;">Name is required.</div>
            </div>
        </div>

        <hr class="mb-4" />
        <button class="btn btn-primary btn-lg btn-block" type="submit">Save</button>
    </form>
</div>`
}

func baseTemplate() string {
	return `{{ define "dashboard" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width">
    <link href="https://unpkg.com/bootstrap@5.0.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-wEmeIV1mKuiNpC+IOBjI7aAzPcEZeedi5yW5f2yOq55WWLwNGmvvx4Um1vskeMj0" crossorigin="anonymous">

    <title>TODO{{ if ne .PageDescription "" }} - {{ .PageDescription }}{{ end }}</title>

    <!-- <meta name="description" content="{{ .PageDescription }}">                 -->
    <!-- <meta property="og:title" content="{{ .PageTitle }}">                      -->
    <!-- <meta property="og:description" content="{{ .PageDescription }}">          -->
    <!-- {{ if ne .PageImagePreview "" }}<meta property="og:image" content="{{ .PageImagePreview }}">{{ end }} -->
    <!-- {{ if and (ne .PageImagePreview "") (ne .PageImagePreviewDescription "") }}<meta property="og:image:alt" content="{{ .PageImagePreviewDescription }}">{{ end }} -->
    <!-- <meta property="og:locale" content="en_GB">                                        -->
    <!-- <meta property="og:type" content="website">                                        -->
    <!-- <meta name="twitter:card" content="summary_large_image">                           -->
    <!-- <meta property="og:url" content="https://www.mywebsite.com/page">                  -->
    <!-- <link rel="canonical" href="https://www.mywebsite.com/page">                       -->

    <link rel="icon" href="/favicon.svg" type="image/svg+xml">
    <link rel="apple-touch-icon" href="/apple-touch-icon.png">
    <!-- <link rel="manifest" href="/my.webmanifest">                                       -->
    <!-- <meta name="theme-color" content="#FF00FF">                                        -->
    <script src="https://js.stripe.com/v3/"></script>
</head>
    <body>
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="/">TODO</a>
                {{ if not .IsLoggedIn }}
                <div class="d-flex">
                    <div class="collapse navbar-collapse" id="navbarNav">
                        <ul class="navbar-nav">
                            <li class="nav-item">
                                <a id="loginLink" class="nav-link" hx-target="#content" hx-push-url="/login" hx-get="/components/login_prompt">{{ translate "callsToAction.signIn" }}</a>
                            </li>
                            <li class="nav-item">
                                <a id="registerLink" class="nav-link" hx-target="#content" hx-push-url="/register" hx-get="/components/registration_prompt">{{ translate "callsToAction.register" }}</a>
                            </li>
                        </ul>
                    </div>
                </div>
                {{ else }}
                <div class="d-flex">
                    <div class="collapse navbar-collapse" id="navbarNav">
                        <ul class="navbar-nav">
                            <li class="nav-item">
                                <a class="nav-link" id="logoutLink" hx-post="/logout">{{ translate "callsToAction.logOut" }}</a>
                            </li>
                        </ul>
                    </div>
                </div>
                {{ end }}
            </div>
        </nav>

        <div class="container-fluid">
            <div class="row">
                <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                    <div class="position-sticky pt-3">
                        <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
                            <span>Things</span>
                        </h6>
                        <ul class="nav flex-column">
                            <li class="nav-item">
                                <a class="nav-link" hx-target="#content" hx-push-url="/items" hx-params="*" hx-get="/dashboard_pages/items">
                                    📃 Items
                                </a>
                            </li>
                            <li class="nav-item">
                                <a class="nav-link"  aria-current="page" hx-target="#content" hx-push-url="/api_clients" hx-params="*" hx-get="/dashboard_pages/api_clients">
                                    🤖 API Clients
                                </a>
                            </li>
                        </ul>
                        <hr>
                        <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
                            <span>Account</span>
                        </h6>
                        <ul class="nav flex-column">
                            <li class="nav-item">
                                <a class="nav-link"  aria-current="page" hx-target="#content" hx-push-url="/account/webhooks" hx-params="*" hx-get="/dashboard_pages/account/webhooks">
                                    🕸️ Webhooks
                                </a>
                            </li>
                            <li class="nav-item">
                                <a class="nav-link" hx-target="#content" hx-push-url="/account/settings" hx-params="*" hx-get="/dashboard_pages/account/settings">
                                    ⚙ Settings
                                </a>
                            </li>
                        </ul>
                        <hr>
                        <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
                            <span>User</span>
                        </h6>
                        <ul class="nav flex-column mb-2">
                            <li class="nav-item">
                                <a class="nav-link" hx-target="#content" hx-push-url="/accounts" hx-params="*" hx-get="/dashboard_pages/accounts">
                                    📚 Accounts
                                </a>
                            </li>
                        </ul>
                        <ul class="nav flex-column mb-2">
                            <li class="nav-item">
                                <a class="nav-link" hx-target="#content" hx-push-url="/user/settings" hx-params="*" hx-get="/dashboard_pages/user/settings">
                                    ⚙ Settings
                                </a>
                            </li>
                        </ul>
                        {{ if .IsServiceAdmin }}
                        <hr>
                        <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
                            <span>Admin</span>
                        </h6>
                        <ul class="nav flex-column mb-2">
                            <li class="nav-item">
                                <a class="nav-link" hx-target="#content" hx-push-url="/admin/users" hx-params="*" hx-get="/dashboard_pages/admin/users">
                                    👥 Users
                                </a>
                            </li>
                        </ul>
                        <ul class="nav flex-column mb-2">
                            <li class="nav-item">
                                <a class="nav-link" hx-target="#content" hx-push-url="/admin/settings" hx-params="*" hx-get="/dashboard_pages/admin/settings">
                                    ⚙ Settings
                                </a>
                            </li>
                        </ul>
                        {{ end }}
                    </div>
                </nav>

                <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                    <div id="content">
                        {{ block "content" .ContentData }}{{ end}}
                    </div>
                </main>
            </div>
        </div>

        <script src="https://unpkg.com/htmx.org@1.3.3" integrity="sha384-QrlPmoLqMVfnV4lzjmvamY0Sv/Am8ca1W7veO++Sp6PiIGixqkD+0xZ955Nc03qO" crossorigin="anonymous"></script>
    </body>
</html>
{{ end }}`
}

func englishTranslationsToml() string {
	return `[testing.translation]
description = "A translation to invoke for tests."
other = ":)"

[callsToAction.signIn]
description = "sign in call to action."
other = "Sign In"

[callsToAction.register]
description = "registration call to action."
other = "Register"

[callsToAction.logOut]
description = "logout call to action."
other = "Log Out"`
}
