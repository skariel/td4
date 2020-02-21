<script>

	import {onMount, getContext} from 'svelte';
	import {post, get} from './utils';

	let user = getContext('user');
	let tests = [];
	let tname = "";

	onMount(initial_load);

	async function initial_load() {
		const res = await get(user, 'tests')
		if (res.status === 200) {
			if (res.data==null) {
				res.data = []
			}
			tests = res.data;
		}
	}

	async function create_test() {
		const res = await post(user, 'create_test', {title:tname, descr:""})
		if (res.status == 200) {
			res.data.avatar = user.avatar
			tests = [res.data, ...tests];
			tname = ""
		}
	}

</script>

<style>
</style>

<svelte:head>
	<title>Sapper project template</title>
</svelte:head>

<input bind:value={tname} placeholder="Test Title">
<button on:click={()=>create_test()} disabled={tname.length == 0}>
	add test
</button>

<h2 style="margin:35px">All Tests</h2>


{#each tests as t }
	<div style="display:flex; align-items:center;margin-top:20px;border:1px solid #333333; border-radius:15px;padding:10px">
		<img style="width:40px; height:40px; margin-right:20px;" src={t.avatar} alt="avatar"/>
		<h3>{t.id}: {t.title}</h3>
		<a style="margin-left:10px;" href={"/test/"+t.id}>more...</a>
	</div>
{/each}
