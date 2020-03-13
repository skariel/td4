<!-- TODO: error handling -->
<!-- TODO: my tests -->
<!-- TODO: my solutions -->
<!-- TODO: delete test / solution -->
<!-- TODO: user profile / stats -->
<!-- TODO: private tests -->
<!-- TODO: draft / publish tests -->

<script>
	import {onMount, getContext} from 'svelte';
	import {get} from '../routes/utils';
	import { goto } from '@sapper/app';
	import TestCard from '../components/TestCard.svelte'

	let user = getContext('user');
	let tests = [];

	const url = new URL(window.location)
	let offset = url.searchParams.get("page")
	if (offset == null) {
		offset = "0";
	}
	offset = parseInt(offset)

	onMount(initial_load);

	function initial_load() {
		get(user, 'alltests/'+offset*10)
			.then((r)=>{tests=r.data;})
	}

	function next() {
		offset += 1;
		initial_load()
	}

	function prev() {
		offset -= 1;
		initial_load()
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
    {#if offset==0}
	    <a href="/?page={offset+1}" on:click={next}>Next Page</a>
    {:else if tests.length < 10}
	    <a href="/?page={offset-1}" on:click={prev}>Prev Page</a>
    {:else}        
	    <a href="/?page={offset-1}" on:click={prev}>Prev Page</a>
        <div class="filler" />
	    <a href="/?page={offset+1}" on:click={next}>Next Page</a>
    {/if}
</div>


