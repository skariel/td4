<script>
	import { get, post, getUser, start_invalidate_cache } from './utils';
	import { onMount } from 'svelte';
	import { writable } from "svelte/store";
	import { goto } from '@sapper/app';

	let user = null;
	let test = [];
	let loading = true;
	let utname = null;
	let utdescr = null;
	let utcode = null;

	function get_test_id_from_url() {
		const url = new URL(location)
		let test_id = url.searchParams.get("id")
		return test_id;
	}

    onMount(()=>{
		let test_id = get_test_id_from_url();
		utname = "";
		utdescr = "";
		utcode = "";

		// utname = writable(localStorage.getItem("utname_"+test_id) || "");
		// utname.subscribe(val => localStorage.setItem("utname_"+test_id, val));
		// utdescr = writable(localStorage.getItem("utdescr_"+test_id) || "");
		// utdescr.subscribe(val => localStorage.setItem("utdescr_"+test_id, val));
		// utcode = writable(localStorage.getItem("utcode_"+test_id) || "");
		// utcode.subscribe(val => localStorage.setItem("utcode_"+test_id, val));

		user = getUser();
		load_data();
    })

	async function load_data(solution_id) {
		loading = true;
        let test_id = get_test_id_from_url();
		let r = await get(user, 'test/'+test_id);
        test=r.data;
        if ($utname == "") {
            set_original_title();
        }
        if ($utdescr == "") {
            set_original_descr();
        }
        if ($utcode == "") {
			set_original_code();
		}
		loading = false;
	}

	function set_original_title() {
		utname.set(test.title)
	}

	function set_original_descr() {
		utdescr.set(test.descr)
	}

	function set_original_code() {
		utcode.set(test.code)
	}

	function set_original_content() {
		set_original_title();
		set_original_descr();
		set_original_code();
	}

	async function update_test() {
		const res = await post(user, 'update_test', {
				id   : test.id,
				title: $utname,
				descr: $utdescr,
				code : $utcode,
			})
		if (res.status == 200) {
			utname.set("")
			utdescr.set("")
			utcode.set("")
			start_invalidate_cache()
			goto('/test?id='+test.id+'&page=0')
		}
	}

	function validate(utname, utdescr, utcode) {
		if (utname.length == 0) {
			return "missing title"
		}
		if (utname.length > 256) {
			return `title is too long: ${utname.length} > 256`
		}
		if (utdescr.length == 0) {
			return "description is missing"
		}
		if (utdescr.length > 2048) {
			return `description is too long: ${utname.length} > 2048`
		}
		if (utcode.length == 0) {
			return "code is missing"
		}
		if (utcode.length > 8192) {
			return `code is too long: ${utname.length} > 8192`
		}
		return ""
	}

</script>

<style>
</style>

<svelte:head>
	<title>Updating test {test.id}</title>
</svelte:head>

{#if loading}
	<h1>loading...</h1>
{:else}

	{#if $utname != test.title || $utdescr != test.descr || $utcode != test.code }
		<button on:click={set_original_content}>restore original test</button>
		<h4 style="color:red;">modified!</h4>
	{/if}

	<h4>Title</h4>
	<input bind:value={$utname}>

	<h4 style="margin-top:20px;">Description</h4>
	<textarea style="width:100%; height:200px" bind:value={$utdescr} />

	<h4 style="margin-top:20px;">Code</h4>
	<textarea style="width:100%; height:200px" bind:value={$utcode} />

	<div>
		{#if validate($utname, $utdescr, $utcode).length!=0}
			<h5 style="margin-top:20px; color:red;">{validate($utname, $utdescr, $utcode)}</h5> 
		{/if}
		<button on:click={update_test} disabled={validate($utname, $utdescr, $utcode).length != 0}>
			update test
		</button>
	</div>
{/if}