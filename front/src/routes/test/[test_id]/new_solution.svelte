<script context="module">
	export function preload({ params }) {
        return {
			test_id: params.test_id,
		}
	}
</script>


<script>
	import { onMount, getContext } from 'svelte';
	import { goto } from '@sapper/app';
	import { get, post } from '../../utils';

	export let test_id;

	const user = getContext('user')
	let test  = [];
	let scode = "";

	onMount(initial_load);

	async function initial_load() {
		get(user, 'test/'+test_id)
			.then((r)=>{test=r.data})
	}

	async function create_solution() {
		const res = await post(user, 'create_solution', {
				test_code_id : parseInt(test_id),
				code       : scode,
			})
		if (res.status == 200) {
			goto('/test/'+test_id+'/0')
		}
	}

</script>

<style>
</style>

<title>New solution</title>

<h4>test code:</h4>
<pre><code>
	{test.code}
</code></pre>

<h4>Code</h4>
<textarea style="width:100%; height:200px" bind:value={scode} />

<button on:click={()=>create_solution()} disabled={scode.length == 0}>
	add solution
</button>

