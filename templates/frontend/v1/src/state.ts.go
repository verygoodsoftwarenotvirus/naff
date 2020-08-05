package frontendsrc

const (
	stateDotTSFile = `export class AuthStatus {
    isAuthenticated: boolean;
    isAdmin: boolean;

    constructor() {
      this.isAuthenticated = false;
      this.isAdmin = false;
    }
}

export interface LoginRequest {
    username: string,
    password: string,
    totpToken: string,
}

export interface UserState {
    authStatus: AuthStatus;
}

export interface AppState {
    user: UserState;
}
`
)

func stateDotTS() string {
	return stateDotTSFile
}
