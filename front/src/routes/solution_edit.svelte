<script>
	import { get, post, getUser } from './utils';
	import { onMount } from 'svelte';
	import { writable } from "svelte/store";
	import { goto } from '@sapper/app';

	let user = null;
    let solution = [];
	const uscode = writable(localStorage.getItem("uscode") || "");
	uscode.subscribe(val => localStorage.setItem("uscode", val));

    onMount(()=>{
		user = getUser();
        load_data();
    })

	async function load_data(solution_id) {
		const url = new URL(location)
        solution_id = url.searchParams.get("id")
		let r = await get(user, 'solution/'+solution_id)
        solution=r.data;
        if ($uscode == "") {
            set_original_code();
        }
	}

	async function update_solution() {
		const res = await post(user, 'update_solution', {
				id   : solution.id,
				code : $uscode,
			})
		if (res.status == 200) {
			scode.set("")
			goto('/test?id='+solution.test_code_id+'&page=0')
		}
	}

	function validate(scode) {
		if (scode.length == 0) {
			return "code is missing";
		}
		if (scode.length > 8192) {
			return `code is too long: ${scode.length} > 8192`;
		}
		return "";
    }

    function set_original_code() {
        $uscode = solution.code
    }

</script>

<style>
</style>

<svelte:head>
	<title>Updating solution {solution.test_code_id}</title>
</svelte:head>

<h1>Updating <a href={"/solution?id="+solution.id}>solution {solution.id}</a> for <a href={"/test?id="+solution.test_code_id}>test {solution.test_code_id}</a></h1>

<h4 style="margin-top:20px;">Code</h4>
<button on:click={set_original_code}>restore solution code</button>
<textarea style="width:100%; height:200px" bind:value={$uscode} />

<div>
	{#if validate($uscode).length > 0}
		<h5 style="margin-top:20px; color:red;">{validate($uscode)}</h5>
	{/if}
	<button on:click={update_solution} disabled={validate($uscode).length != 0}>
		update solution
	</button>
</div>
