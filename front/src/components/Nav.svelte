<script>
	import { staticpath, getUser } from '../routes/utils'

	import { onMount } from 'svelte';

	let user = {};
    onMount(()=>{
        user = getUser();
    })

	export let segment;
</script>

<style>
	nav {
		border-bottom: 1px solid rgba(255,62,0,0.1);
		font-weight: 300;
		padding: 0 1em;
		max-width: 60em;
		margin: auto;
	}

	ul {
		margin: 0;
		padding: 0;
		display: flex;
	}

	/* clearfix */
	ul::after {
		content: '';
		display: block;
		clear: both;
	}

	li {
		display: block;
		float: left;
	}

	.selected {
		position: relative;
		display: inline-block;
	}

	.selected::after {
		position: absolute;
		content: '';
		width: calc(100% - 1em);
		height: 2px;
		background-color: rgb(255,62,0);
		display: block;
		bottom: -1px;
	}

	a {
		text-decoration: none;
		padding: 1em 0.5em;
		display: block;
	}

	li.avatar {
		margin-left: auto;
	}

	li.avatar a {
		padding: 0;
	}

	img.avatar {
		height: 2em;
		margin-top: 0.7em;
	}

</style>

<nav>
	<ul>
		<li><a class:selected='{segment === undefined}' href='.'>home</a></li>
		<li><a class:selected='{segment === "about"}' href='about'>about</a></li>
		<li class="avatar">
			{#if user['avatar'] == null}
				<a href="https://localhost:8081/auth/github">
					<img class="avatar" alt="login with github" src={staticpath("login_with_github.png")} />
				</a>
			{:else}
					<img class="avatar" alt="avatar" src={user.avatar} />
			{/if}
		</li>
	</ul>
</nav>
