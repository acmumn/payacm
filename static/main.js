$(() => {
	let amount = 0;
	const handler = StripeCheckout.configure({
	   	key: "pk_test_Mxl12rTX2BQpuAVWg6eEElFp",
		locale: "auto",
		token: onSubmit,
	});

	$("#amount").on("input", onAmountChange);
	$("#pay").click(onPay);
	if(window.location.hash) {
		$("#amount").val(window.location.hash.slice(1));
		onAmountChange();
	}

	function onAmountChange() {
		const value = $("#amount").val();
		if(value.match(/^[0-9]+(\.[0-9]{2})$/)) {
			amount = parseInt(value.replace(".", ""));

			$("#amount-group").addClass("has-success");
			$("#amount-group").removeClass("has-error");
			$("#pay").attr("disabled", false);
		} else {
			$("#amount-group").addClass("has-error");
			$("#amount-group").removeClass("has-success");
			$("#pay").attr("disabled", true);
		}
	}
	function onPay(e) {
		handler.open({
			name: "ACM UMN",
			description: "payacm",
			amount: amount,
		});
		e.preventDefault();
	}
	function onSubmit(token) {
		const body = {
			amount: amount,
			email: token.email,
			token: token.id,
		};
		fetch("/", {
			body: JSON.stringify(body),
			method: "POST"
		}).then(res => {
			return res.json().then(json => {
				if(res.ok)
					return json;
				else
					throw json;
			});
		}).then(res => {
			alert("Success!");
		}).catch(err => {
			console.error(err);
			alert("An error occurred.");
		});
	}
});
