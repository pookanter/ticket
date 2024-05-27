import { httpService } from './http-service';

export namespace IAuthenService {
	export interface Authorization {
		access_token: string;
		refresh_token: string;
	}
}

export const AuthenService = {
	getAuthorization() {
		const token = localStorage.getItem('x-authorization');

		if (token) {
			return JSON.parse(token) as IAuthenService.Authorization;
		}

		return null;
	},
	setAuthorization(token: IAuthenService.Authorization) {
		localStorage.setItem('x-authorization', JSON.stringify(token));
	},

	removeAuthorization() {
		localStorage.removeItem('x-authorization');
	},

	getAccessToken() {
		const token = AuthenService.getAuthorization();

		return token?.access_token;
	},

	getRefreshToken() {
		const token = AuthenService.getAuthorization();

		return token?.refresh_token;
	},

	refreshToken(refreshToken: string) {
		return httpService().post<IAuthenService.Authorization>('/authen-service/refresh-token', {
			refreshToken
		});
	},

	signIn(email: string, password: string) {
		return httpService().post<IAuthenService.Authorization>('/authen-service/sign-in', {
			email,
			password
		});
	},

	signUp(email: string, password: string) {
		return httpService().post<{ message: string; user_id: number }>('/authen-service/sign-up', {
			email,
			password
		});
	}
};
