<script>
	import { get, post, getUser, start_invalidate_cache } from './utils';
	import { onMount } from 'svelte';
	import { writable } from "svelte/store";
	import { goto } from '@sapper/app';

	let user = null;
    let solution = [];
	let loading = true;
	let uscode = null;

	function get_solution_id_from_url() {
		const url = new URL(location)
		let solution_id = url.searchParams.get("id")
		return solution_id;
	}

	onMount(()=>{
		let solution_id = get_solution_id_from_url();
		uscode = writable(localStorage.getItem("uscode"+solution_id) || "");
		uscode.subscribe(val => localStorage.setItem("uscode"+solution_id, val));

		user = getUser();
        load_data();
	})

	async function load_data() {
		loading = true;
        let solution_id = get_solution_id_from_url();
		let r = await get(user, 'solution/'+solution_id)
        solution=r.data;
        if ($uscode == "") {
            set_original_code();
		}
		loading = false;
	}

	function set_original_code() {
		uscode.set(solution.code)
	}

	async function update_solution() {
		const res = await post(user, 'update_solution', {
				id   : solution.id,
				code : $uscode,
			})
		if (res.status == 200) {
			uscode.set("")
			start_invalidate_cache()
			goto('/test?id='+solution.test_code_id+'&page=0')
		}
	}

	function validate(uscode) {
		if (uscode.length == 0) {
			return "code is missing";
		}
		if (uscode.length > 8192) {
			return `code is too long: ${uscode.length} > 8192`;
		}
		return "";
    }


</script>

<style>
</style>

<svelte:head>
	<title>Updating solution {solution.test_code_id}</title>
</svelte:head>

{#if loading}
	<h1>loading...</h1>
{:else}
	<h1>Updating <a href={"/solution?id="+solution.id}>solution {solution.id}</a> for <a href={"/test?id="+solution.test_code_id}>test {solution.test_code_id}</a></h1>

	<h4 style="margin-top:20px;">Code</h4>

	{#if $uscode != solution.code}
		<button on:click={set_original_code}>restore solution code</button>
		<h4 style="color:red;">modified!</h4>
	{/if}
	<textarea style="width:100%; height:200px" bind:value={$uscode} />

	<div>
		{#if validate($uscode).length > 0}
			<h5 style="margin-top:20px; color:red;">{validate($uscode)}</h5>
		{/if}
		<button on:click={update_solution} disabled={validate($uscode).length != 0}>
			update solution
		</button>
	</div>
{/if}
