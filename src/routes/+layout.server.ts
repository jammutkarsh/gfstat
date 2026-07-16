export function load({ cookies }) {
	return { loggedIn: !!cookies.get('gh_token') };
}
