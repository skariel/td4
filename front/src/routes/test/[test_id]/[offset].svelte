<script context="module">
	export function preload({ params }) {
        return {
			test_id: params.test_id,
			offset  : params.offset,
		}
	}
</script>

<script>
	import { onMount, getContext } from 'svelte';
	import { get } from '../../utils';

    export let test_id;
    export let offset;
	let user  = getContext('user');
	let test  = [];
	let solutions = [];

	onMount(initial_load);

	async function initial_load() {
		const _tes = get(user, 'test/'+test_id)
		const _sol = get(user, 'solutions_by_test/'+test_id+'/'+offset)
		test = (await _tes).data
		solutions = (await _sol).data
	}
</script>

<style>
</style>

<svelte:head>
	<title>just a te
	st {test_id}!</title>
</svelte:head>

<a href={'/test/'+test_id+'/new_solution'}>add solution</a>

<h4>test code:</h4>
<pre><code>
	{test.code}
</code></pre>

<h4>solutions:</h4>

{#each solutions as s }
	<div  style="margin:15px;">
		<h5>s:</h5>
		<pre>
			<code>
				{s.code}
			</code>
		</pre>
		<a href="/test/{test_id}/code/{s.id}">code id: {s.id}</a>
	</div>
{/each}
