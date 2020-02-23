<script context="module">
	export function preload({ params }) {
        return {
            test_id: params.test_id,
            code_id: params.code_id
        }
	}
</script>

<script>
    export let test_id;
    export let code_id;
	let user = getContext('user');
    let code = ""
    let solutions = []
	let new_solution = ""

	import {onMount, getContext} from 'svelte';
	import {post, get} from '../../../utils';

	onMount(initial_load);

	async function initial_load() {
		const _cod = get(user, 'code/'+code_id)
		const _sol = get(user, 'code/'+code_id+'/solutions')
		const cod = await _cod
		const sol = await _sol
		if (cod.status === 200) {
			code = cod.data;
		}
		if (sol.status === 200) {
			solutions = sol.data;
			if (solutions == null) {
				solutions = []
			}
        }
        console.log(code)
        console.log(solutions)
	}

	async function create_solution() {
		const res = await post(user, 'create_solution', {test_code_id:parseInt(code_id), code:new_solution})
		if (res.status == 200) {
			solutions = [res.data, ...solutions];
			new_solution = ""
		}
	}


</script>

<style>
</style>

<svelte:head>
	<title>just a COdeID {code_id}!</title>
</svelte:head>

<h3>Test {test_id} > Test Code (id={code_id}):</h3>
<div>
    <code>
        {code.code}
    </code>
</div>

<h3>Solutions:</h3>

<button on:click={()=>create_solution()} disabled={new_solution.length == 0}>
	add solution!
</button>
<div style="margin:15px;">
	<textarea style="width:100%; height:200px" bind:value={new_solution} />
</div>


{#each solutions as s}
    <div>
        <code>
            {s.code}
        </code>
    </div>
{/each}