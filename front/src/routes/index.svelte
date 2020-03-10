<!-- TODO: error handling -->
<!-- TODO: my tests -->
<!-- TODO: my solutions -->
<!-- TODO: delete test / solution -->
<!-- TODO: user profile / stats -->
<!-- TODO: private tests -->
<!-- TODO: draft / publish tests -->

<script>

	import {onMount, getContext} from 'svelte';
	import {post, get} from './utils';
	import TestCard from '../components/TestCard.svelte'

	let user = getContext('user');
	let tests = [];
	let offset = 0;

	onMount(initial_load);

	async function initial_load() {
		const res = await get(user, 'alltests/'+offset)
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