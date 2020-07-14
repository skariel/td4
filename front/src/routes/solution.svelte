<script>
	import { onMount } from 'svelte';
	import { get, getUser } from './utils';
	import { goto } from '@sapper/app';

    let user          = {};
    let solution_id   = 0;
	let solution      = [];
	let results       = [];
	let loading       = 2;

	onMount(async ()=>{
        user = getUser();
		load_data()
	});

	async function load_data() {
		loading = 2;
		const url = new URL(location)
        solution_id = url.searchParams.get("id")
		let r = await get(user, 'solution/'+solution_id)
		solution=r.data;
		r = await get(user, 'results_by_run/'+solution.run_id)
		results=r.data;
		loading = 0;
	}
</script>

<style>
    .top {
        display:       flex;
    }
    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;
    }

    .solutionid {
        margin-left: auto;
    }

	.code {
		background-color: #333333;
		min-height: 100px;
	}

	code {
		background-color: #00000000;
		color: antiquewhite;
	}
	.title {
		display: flex;
		align-items: center;
	}
</style>

<svelte:head>
	<title>Solution {solution_id}</title>
</svelte:head>

{#if loading>0}
	<h1>loading...</h1>
{:else}
	<div class="top">
		<img class="avatar" src={solution.avatar} alt="avatar"/>
		<h4 class="displayname">{solution.display_name}</h4>
		<h4 class="solutionid">solution {solution.id}</h4>
	</div>

	<h4 style="margin-top:10px;">Solution for <a href={"/test?id="+solution.test_code_id}>test {solution.test_code_id}</a></h4>

	<h4>status: {solution.status} {#if solution.status=='stop'} (timeout) {/if}</h4>

	{#if user.display_name==solution.display_name}
		<button on:click={goto("/solution_edit?id="+solution.id)}>Edit Code</button>
	{/if}

	<pre class="code">
		<code>
			{solution.code}
		</code>
	</pre>

	{#if results.length > 0}
		<div class="title">
			<h1>Results:</h1>
		</div>
		<pre class="code">
			<code>
				{results[0].output.String}
			</code>
		</pre>
	{:else}
		<div class="title">
			<h1>No result yet</h1>
		</div>
	{/if}

{/if}