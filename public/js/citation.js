/**
 * Created by fv on 13.03.17.
 */

$(document).ready(function() {
    $('select').material_select();

    $('.datepicker').pickadate({
        selectMonths: true, // Creates a dropdown to control month
        selectYears: 15 // Creates a dropdown of 15 years to control year
    });

    CitationDB.Bootstrap();
});

var CitationDB = (function ()
{

    var app         = {};

    app.UrlPath         = "/";
    app.UrlPathParts    = [];

    app.Bootstrap = function () {
        app.UrlPath         = window.location.pathname;
        app.UrlPathParts    = app.UrlPath.split( '/' );

        app.UrlPathParts    = app.UrlPathParts.slice(1);

        switch (app.UrlPathParts[0])
        {
            case "books":
                app.Books.Init();
                break;
            case "quotes":
                app.Quotes.Init();
                break;
        }

    };

    return app;
})();

CitationDB.Books = (function ()
{
    var _self = {};

    _self.API = (function() {
        var _api = {};

        _api.searchForBook = function (name, numOfItems) {
            var def = $.Deferred();

            if (!numOfItems)
                numOfItems = 10;

            $.post("/books/search", { "query": name, "limit": numOfItems}).success(function (result) {

                console.debug(result);
                def.resolve(result);
            }).fail(function (msg) {
                console.debug(msg);
                def.reject(msg);
            });

            return def.promise();
        };

        return _api;
    })();

    _self.Init = function ()
    {
        if (CitationDB.UrlPathParts.length > 1)
        {
            switch (CitationDB.UrlPathParts[1])
            {
                case "list":
                    _self.BuildCollectionView();
                    break;
            }
        }
        else
        {
            // Invalid path
        }
    };

    _self.BuildCollectionView = function ()
    {
        $('.cdb-books-list-entry').on('click', _self.onBooksListEntryClicked);
    };

    _self.onBooksListEntryClicked = function (ev)
    {
        var entry = $(ev.currentTarget);
        var id = entry.attr('data-id');

        window.location.href = "/books/edit/" + id;
    };

    return _self;
})();

CitationDB.Quotes = (function () {
    var _self = {};

    _self.Init = function () {
        if (CitationDB.UrlPathParts.length > 1)
        {
            switch (CitationDB.UrlPathParts[1])
            {
                case "add":
                    _self.buildAddView();
                    break;
            }
        }
        else
        {
            // Invalid path
        }
    };

    _self.buildAddView = function () {
        $('#quoteBook').on("input", _self.onBookInputContentChanged);
    };

    _self.buildAutoCompleteView = function (bookArr) {

        if (!bookArr)
            return;

        var data = {};
        var bookElemC = bookArr.length;

        for (var i=0; i<bookElemC; i++)
        {
            data[bookArr[i].Title] = bookArr[i].Id;
        }

        $('#quoteBook').autocomplete({
            data: data
        });

        $('.autocomplete-content').on("click", function (ev) {
            console.debug(ev.target);
            var id = $(ev.target).parents('li').find('img').attr('src');
            console.debug(id);
            $('#selectedBook').attr('value', id);
        })
    };

    _self.onBookInputContentChanged = function (ev) {
        var val = $(ev.currentTarget).val();

        CitationDB.Books.API.searchForBook(val).done(function (response) {
            _self.buildAutoCompleteView(response.results);
        });
    };

    return _self;
})();
