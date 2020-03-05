<script>

	import {onMount, getContext} from 'svelte';
	import {post, get} from './utils';

	let user = getContext('user');
	let tests = [];
	let offset = 0;

	onMount(initial_load);

	async function initial_load() {
		const res = await get(user, 'alltests/'+offset)
		tests = res.data;
		console.log(tests)
	}

</script>

<style>
</style>

<svelte:head>
	<title>Tesoto</title>
</svelte:head>


<a href="/test/new">New Test</a>

<h2 style="margin:35px">All Tests</h2>
{#each tests as t }
	<div style="display:flex; align-items:center;margin-top:20px;border:1px solid #333333; border-radius:15px;padding:10px">
		<img style="width:40px; height:40px; margin-right:20px;" src={t.avatar} alt="avatar"/>
		<h3>{t.id}: {t.title}</h3>
		<a style="margin-left:10px;" href={"/test/"+t.id}>more...</a>
	</div>
{/each}
