<script>
	import { beforeUpdate, setContext } from 'svelte'
	beforeUpdate(updateUser)

	import Nav from '../components/Nav.svelte';

	export let segment;
	

	function updateUser() {
		let user = {};
		let cookies = document.cookie.split(';')
		for (const ix in cookies) {
			const c = cookies[ix]
			const cs = c.split('=')
			user[cs[0].trim().replace('user_', '')] = cs[1];
		}
		setContext('user', user)
	}
</script>

<style>
	main {
		position: relative;
		max-width: 56em;
		background-color: white;
		padding: 2em;
		margin: 0 auto;
		box-sizing: border-box;
	}
</style>

<Nav segment={segment}/>

<main>
	<slot></slot>
</main>