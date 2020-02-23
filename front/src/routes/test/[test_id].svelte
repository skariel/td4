<script context="module">
	export function preload({ params }) {
        return {test_id: params.test_id}
	}
</script>

<script>
    export let test_id;
	let user  = getContext('user');
	let test  = {title:""}
	let codes = []
	let new_code = ""

	import {onMount, getContext} from 'svelte';
	import {post, get} from '../utils';

	onMount(initial_load);

	async function initial_load() {
		const _tes = get(user, 'test/'+test_id)
		const _cod = get(user, 'test/'+test_id+'/codes')
		const tes = await _tes
		const cod = await _cod
		if (tes.status === 200) {
			test = tes.data;
		}
		if (cod.status === 200) {
			codes = cod.data;
			if (codes == null) {
				codes = []
			}
		}
	}


	async function create_code() {
		const res = await post(user, 'create_code', {test_id:parseInt(test_id), code:new_code})
		if (res.status == 200) {
			codes = [res.data, ...codes];
			new_code = ""
		}
	}


</script>

<style>
</style>

<svelte:head>
	<title>just a test {test_id}!</title>
</svelte:head>

<h3>{test_id}: {test.title}</h3>
<button on:click={()=>create_code()} disabled={new_code.length == 0}>
	add code!
</button>
<div  style="margin:15px;">
	<textarea bind:value={new_code} />
</div>

{#each codes as c }
	<div  style="margin:15px;">
		<code>
			{c.code}
		</code>
		<a href="/test/{test_id}/code/{c.id}">code id: {c.id}</a>
	</div>
{/each}
