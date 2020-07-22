<script>
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '@sapper/app';

	import { get, del, getUser, loginpath, start_invalidate_cache } from './utils';
	import SolutionCard from '../components/SolutionCard.svelte'

	let user      = {};
    let page      = 0;
    let test_id   = 0;
	let test      = [];
	let solutions = [];
	let loading   = 2;

	onMount(async ()=>{
        user = getUser();
		window.addEventListener("locationchange", load_data);
		load_data()
	});

	onDestroy(async ()=>{
        // TODO: otherwise running on node?! a bug?!
        if (process.browser == true) {
            window.removeEventListener("locationchange", load_data);
        }
	})

	function params(page) {
		return `/test?id=${test_id}&page=${page}`
	}

	function load_data() {
		loading = 2;
		const pname_wo_trailin_slash = window.location.pathname.replace(/\/+$/, '');
        if (pname_wo_trailin_slash != '/test') {
            return
        }
		const url = new URL(location)
		page = url.searchParams.get("page")
		if (page == null) {
			page = "0";
		}
		page = parseInt(page)
        test_id = url.searchParams.get("id")
		get(user, 'test/'+test_id)
			.then((r)=>{test=r.data; loading -= 1;})
		get(user, 'solutions_by_test/'+test_id+'/'+page*10)
			.then((r)=>{solutions=r.data; loading -= 1;})
	}

	async function delete_test() {
		const res = await del(user, 'delete_test/'+test.id);
		if (res.status == 200) {
			start_invalidate_cache();
			goto('/')
		}
	}

</script>

<style>
	.solutions {
 		column-count: 2;
		column-gap: 1em;
	}

	.solutioncard {
		display: inline-block;
		margin-top: 10px;
		width: 100%;
	}

    .top {
        display:       flex;
    }
    .teststat {
        display:       flex;
		margin-top:    0px; 
    }
    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;
    }

    .testid {
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
	.title a {
		margin-left: auto;
	}
	.bottom {
		display: flex;
		justify-content: center;
		margin-top: 20px;
	}
    .filler {
        width: 20px;
    }


</style>

<svelte:head>
	<title>Test {test_id}</title>
</svelte:head>

{#if loading>0}
	<h1>Loading...</h1>
{:else}
	<div class="top">
		<img class="avatar" src={test.avatar} alt="avatar"/>
		<h4 class="displayname">{test.display_name}</h4>
		<h4 class="testid">test {test.id}</h4>
	</div>

	<h3>Title: {test.title}</h3>
	<p>Description: {test.descr}</p>

	{#if user.display_name==test.display_name}
		<button on:click={goto("/test_edit?id="+test.id)}>Edit Test</button>
		<button on:click={delete_test}>Permanently delete test</button>
	{/if}

	<pre class="code">
		<code>
			{test.code}
		</code>
	</pre>

	<div class="title">
		{#if solutions.length > 0}
			<h1>All solutions</h1>
		{:else}
			{#if page == 0}
				<h1>No solutions yet!</h1>
			{:else}
				<h1>No solutions in this page!</h1>
			{/if}
		{/if}
		{#if user['avatar'] != null}
			<a href={"/new_solution?test_id="+test_id}>Add Solution</a>
		{:else}
			<a href={loginpath()}>Login to add a solution</a>
		{/if}
	</div>

	{#if solutions.length > 0}
		<div class="teststat">
			<h4>fail: {test.total_fail}</h4>
			<h4 style="margin-left:10px;">pass: {test.total_pass}</h4>
			<h4 style="margin-left:10px;">pending: {test.total_pending}</h4>
			<h4 style="margin-left:10px;">wip: {test.total_wip}</h4>
		</div>
	{/if}

	<div class="solutions">
		{#each solutions as s }
			<div class="solutioncard">
				<SolutionCard solution={s} />
			</div>
		{/each}
	</div>

	<div class="bottom">
		{#if page==0}
			{#if solutions.length == 10}
				<a href="{params(page+1)}">Next Page</a>
			{/if}
		{:else if solutions.length < 10}
			<a href="{params(page-1)}">Prev Page</a>
		{:else}
			<a href="{params(page-1)}">Prev Page</a>
			<div class="filler" />
			<a href="{params(page+1)}">Next Page</a>
		{/if}
	</div>


{/if}