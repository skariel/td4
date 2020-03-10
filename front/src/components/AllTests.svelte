<!-- TODO: loading message -->

<script>

	import {onMount, getContext} from 'svelte';
	import {post, get} from '../routes/utils';
	import { goto } from '@sapper/app';
	import TestCard from './TestCard.svelte'

	let user = getContext('user');
	let tests = [];
	export let offset = 0;

	onMount(initial_load);

	async function initial_load() {
		const res = await get(user, 'alltests/'+offset*10)
		tests = res.data;
	}

</script>

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

<!-- TODO: file bug. The href here is necessary otherwise onmount will not be always called! -->
<div class="bottom">
    {#if offset=="0"}
	    <a href="/{parseInt(offset)+1}" onclick="location.href='/{parseInt(offset)+1}'">Next Page</a>
    {:else if tests.length < 10}
	    <a href="/{parseInt(offset)-1}" onclick="location.href='/{parseInt(offset)-1}'">Prev Page</a>
    {:else}        
	    <a href="/{parseInt(offset)-1}" onclick="location.href='/{parseInt(offset)-1}'">Prev Page</a>
        <div class="filler" />
	    <a href="/{parseInt(offset)+1}" onclick="location.href='/{parseInt(offset)+1}'">Next Page</a>
    {/if}
</div>
