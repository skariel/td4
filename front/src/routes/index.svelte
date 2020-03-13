<!-- TODO: error handling -->
<!-- TODO: my tests -->
<!-- TODO: my solutions -->
<!-- TODO: delete test / solution -->
<!-- TODO: user profile / stats -->
<!-- TODO: private tests -->
<!-- TODO: draft / publish tests -->

<script>
	import { onMount, getContext } from 'svelte';
	import { get, init_location_change_event } from '../routes/utils';
	import TestCard from '../components/TestCard.svelte'

	let user = getContext('user');
	let tests = [];

	let page = null;
	let href = null;

	onMount(()=>{
		init_location_change_event()
		window.addEventListener('locationchange', load_data)
		load_data()
	});

	function load_data() {
		const url = new URL(location)
		page = url.searchParams.get("page")
		if (page == null) {
			page = "0";
		}
		page = parseInt(page)
		get(user, 'alltests/'+page*10)
			.then((r)=>{tests=r.data;})
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
	<title>TesTus</title>
</svelte:head>

<div class="title">
	<h1>All Tests</h1>
	<a href="/test/new">New Test</a>
</div>

<div class="tests">
	{#each tests as t }
		<div class="testcard">
			<TestCard test={t} />
		</div>
	{/each}
</div>

<div class="bottom">
    {#if page==0}
	    <a href="/?page={page+1}">Next Page</a>
    {:else if tests.length < 10}
	    <a href="/?page={page-1}">Prev Page</a>
    {:else}        
	    <a href="/?page={page-1}">Prev Page</a>
        <div class="filler" />
	    <a href="/?page={page+1}">Next Page</a>
    {/if}
</div>