<script>
	import { onMount, onDestroy } from 'svelte';
	import { get, getUser, loginpath } from './utils';
	import TestCard from '../components/TestCard.svelte'

	let tests = [];
	let user = {};
	let page = null;
	let loading = true;
	let user_filter = null;

	onMount(()=>{
		user = getUser();
		window.addEventListener("locationchange", load_data);
		load_data()
	});

	onDestroy(()=>{
        // TODO: otherwise running on node?! a bug?!
        if (process.browser == true) {
			window.removeEventListener("locationchange", load_data);
		}
	})

	function load_data() {
		loading = true;
        if (window.location.pathname != '/') {
            return
        }
		const url = new URL(location)
		page = url.searchParams.get("page")
		if (page == null) {
			page = "0";
		}
		page = parseInt(page)
		var href;
		user_filter = url.searchParams.get("user")
		if (user_filter==null) {
			href = 'alltests/'+page*10;
		}
		else {
			href = 'alltests_by_user/'+page*10+'/'+user_filter;
		}
		get(user, href)
			.then((r)=>{tests=r.data; loading = false;})
	}
</script>

<title>TesTus</title>

<style>
	.tests {
 		column-count: 2;
		column-gap: 1em;
	}

	.title {
		display: flex;
		align-items: center;
	}

	.title a {
		margin-left: auto;
	}

	.testcard {
		display: inline-block;
		margin-top: 10px;
		width: 100%;
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
	<title>Tests</title>
</svelte:head>

<div style="visibility: hidden;">
<!-- exporting with sapper -->
	<a href="/test">exporting</a>
	<a href="/solution">exporting</a>
	<a href="/new_solution">exporting</a>
	<a href="/new_test">exporting</a>
	<a href="/solution_edit">exporting</a>
	<a href="/test_edit">exporting</a>
</div>

<div class="title">
	{#if user_filter == null}
		<h1>Showing all tests</h1>
	{:else}
		<h1>Showing tests for {user_filter}</h1>
	{/if}
	{#if user['avatar'] != null}
		{#if user_filter != null}
			<a style="margin-right?:15px; margin-left:auto;" href="/?page=0" on:click={()=>{load_data();}}>all tests</a>
		{:else}
			<a style="margin-right?:15px; margin-left:auto;" href="/?page=0&user={user.display_name}">my tests</a>
		{/if}
		<a style="margin-right?:15px; margin-left:15px;" href="/new_test">Add Test</a>
	{:else}
		<a href={loginpath()}>Login to add a test</a>
	{/if}
</div>

{#if loading}
	<h2>Loading...</h2>
{:else}
	{#if tests.length == 0}
		{#if page == 0}
			<h2>None yet!</h2>
		{:else}
			<h2>None on this page!</h2>
		{/if}
	{/if}
{/if}

<div class="tests">
	{#each tests as t }
		<div class="testcard">
			<TestCard test={t} />
		</div>
	{/each}
</div>

<div class="bottom">
    {#if page==0}
		{#if tests.length == 10}
	    	<a href="/?page={page+1}">Next Page</a>
		{/if}
    {:else if tests.length < 10}
	    <a href="/?page={page-1}">Prev Page</a>
    {:else}
	    <a href="/?page={page-1}">Prev Page</a>
        <div class="filler" />
	    <a href="/?page={page+1}">Next Page</a>
    {/if}
</div>

