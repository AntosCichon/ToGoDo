$(".add").on('click', function() {
    let title = $("#todo-title").val()
    let color = $("#todo-color").val()
    $.ajax({
        url: "/add",
        type: "POST",
        data: {
            title: title,
            color: color,
        },
        success: function(response) {
            $("#todos").append(`
                <div class="item" data-id="${response.id}" data-modifier="${response.modifier}" data-color="${response.color}">
                    <p>${response.title}</p>
                    <button class="remove">X</button>
                </div>
            `);
            $("#todo-title").val("");
            colorize();
        },
        error: function(response) {
        }
    });
});

$(document).on('click', '.remove', function() {
    let id = $(this).parent().attr('data-id')
        $.ajax({
            url: "/remove",
            type: "POST",
            data: {
                id: id,
            },
            success: function(response) {
                let html = "";
                response.forEach(entry => {
                    html += `
                    <div class="item" data-id="${entry.id}" data-modifier="${entry.modifier}" data-color="${entry.color}">
                        <p>${entry.title}</p>
                        <button class="remove">X</button>
                    </div>
                `});
                $("#todos").html(html);
                colorize();
            },
            error: function(response) {
            }
        });
});

$(document).ready(function() {
    colorize();
    $('#todo-color').on('input', function() {
      const hue = $(this).val();
      const hsl = `hsl(${hue}, 60%, 60%)`;
      const hslBackground = `hsl(${hue}, 20%, 30%)`;
      $("#input button").css({
        'color': hsl,
        'border': `2px solid ${hsl}`,
        'background-color': hslBackground,
    });
    });
  });

$(document).on('click', '.item p', function() {
    let entry = $(this).parent();
    let entryId = entry.attr("data-id");
    let modifier = entry.attr("data-modifier");
    newModifier = modifier * -1 + 1;
    entry.attr("data-modifier", newModifier);
    $.ajax({
        url: "/modify",
        type: "POST",
        data: {
            id: entryId,
            modifier: newModifier,
        },
        success: function(response) {
            let html = "";
            response.forEach(entry => {
                console.log(html)
                html += `
                <div class="item" data-id="${entry.id}" data-modifier="${entry.modifier}" data-color="${entry.color}">
                    <p>${entry.title}</p>
                    <button class="remove">X</button>
                </div>
            `});
            $("#todos").html(html);
            colorize();
        },
        error: function(response) {
        }
    });
})

function colorize() {
    $(".item").each(function() {
        const hue = $(this).attr("data-color");
        const hsl = `hsl(${hue}, 60%, 60%)`;
        const hslBackground = `hsl(${hue}, 20%, 30%)`;
        $(this).css({
            'border': `2px solid ${hsl}`,
            'background-color': hslBackground,
        });
        $(this).children().each(function() {
            $(this).css('color', hsl);
        })
    })
}