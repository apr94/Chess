
$(document).ready(function() {
	$('.StartGameButton').click(function() {
		document.getElementById("StartGameForm").submit();
	});
});

$(document).ready(function() {
	$('.NewGameButton').click(function() {
		document.getElementById("NewGameForm").submit();
	});
});






$(document).ready(function ()
{

	for(var i = 0; i < 61; i++){
		$('#MinutesList').append("<option value = \"" + i + "\">" + i + "</option>");
		$('#SecondsList').append("<option value = \"" + i + "\">" + i + "</option>");
		$('#IncrementList').append("<option value = \"" + i + "\">" + i + "</option>");	
	}



	$('#TimeControls').hide();
})


$(document).ready(function ()
{
	$("input[name='Timing']").change(radioValueChanged);
})

function radioValueChanged()
{
	radioValue = $(this).val();

	if($(this).is(":checked") && radioValue == "UnTimed")
	{
		$('#TimeControls').hide();
	}
	else
	{
		$('#TimeControls').show();
	}
} 
