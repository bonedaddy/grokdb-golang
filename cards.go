package main

import (
    "database/sql"
    "errors"
    "net/http"
    "strconv"
    "strings"

    // 3rd-party
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

/* variables */

// errors
var ErrCardNoSuchCard = errors.New("cards: no such card of given id")
var ErrCardNoCardsByDeck = errors.New("cards: deck has no cards")
var ErrCardPageOutOfBounds = errors.New("cards: page is out of bounds")

/* types */

type CardProps struct {
    Title       string
    Description string
    Front       string
    Back        string
    Deck        uint
}

type CardRow struct {
    ID          uint `db:"card_id"`
    Title       string
    Description string
    Front       string
    Back        string
    Deck        uint  `db:"deck"`
    CreatedAt   int64 `db:"created_at"`
    UpdatedAt   int64 `db:"updated_at"`
}

type CardPOSTRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
    Front       string `json:"front"`
    Back        string `json:"back"`
    Deck        uint   `json:"deck" binding:"required,min=1"`
}

/* REST Handlers */

// GET /cards/:id
//
// Params:
// id: a unique, positive integer that is the identifier of the assocoated deck
func CardGET(db *sqlx.DB, ctx *gin.Context) {

    var err error

    // parse and validate id param
    var cardIDString string = strings.ToLower(ctx.Param("id"))

    _cardID, err := strconv.ParseUint(cardIDString, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": err.Error(),
            "userMessage":      "given id is invalid",
        })
        ctx.Error(err)
        return
    }
    var cardID uint = uint(_cardID)

    // fetch card row from the db

    var fetchedCardRow *CardRow
    fetchedCardRow, err = GetCard(db, cardID)
    switch {
    case err == ErrCardNoSuchCard:
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":           http.StatusNotFound,
            "developerMessage": err.Error(),
            "userMessage":      "cannot find card by id",
        })
        ctx.Error(err)
        return
    case err != nil:
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card",
        })
        ctx.Error(err)
        return
    }

    // fetch card's score
    var fetchedCardScore *CardScoreRow
    fetchedCardScore, err = GetCardScoreRecord(db, cardID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card score record",
        })
        ctx.Error(err)
        return
    }

    var fetchedStashes []uint
    fetchedStashes, err = StashesByCard(db, cardID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card stashes",
        })
        ctx.Error(err)
        return
    }

    var cardrow gin.H = CardRowToResponse(db, fetchedCardRow)
    var cardscore gin.H = CardScoreToResponse(fetchedCardScore)

    ctx.JSON(http.StatusOK, MergeResponses(
        &cardrow,
        &gin.H{"review": cardscore},
        &gin.H{"stashes": fetchedStashes},
    ))
}

// DELETE /cards/:id
//
func CardDELETE(db *sqlx.DB, ctx *gin.Context) {

    var err error

    // parse and validate id param
    var cardIDString string = strings.ToLower(ctx.Param("id"))

    _cardID, err := strconv.ParseUint(cardIDString, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": err.Error(),
            "userMessage":      "given id is invalid",
        })
        ctx.Error(err)
        return
    }
    var cardID uint = uint(_cardID)

    // fetch card row from the db

    _, err = GetCard(db, cardID)
    switch {
    case err == ErrCardNoSuchCard:
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":           http.StatusNotFound,
            "developerMessage": err.Error(),
            "userMessage":      "cannot find card by id",
        })
        ctx.Error(err)
        return
    case err != nil:
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card",
        })
        ctx.Error(err)
        return
    }

    // delete card
    err = DeleteCard(db, cardID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to delete card",
        })
        ctx.Error(err)
        return
    }

    // ensure card is deleted

    _, err = GetCard(db, cardID)
    switch {
    case err == ErrCardNoSuchCard:
        // success
        ctx.Writer.WriteHeader(http.StatusNoContent)
        return
    case err != nil:
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to check if card is deleted",
        })
        ctx.Error(err)
        return
    }

    ctx.JSON(http.StatusInternalServerError, gin.H{
        "status":           http.StatusInternalServerError,
        "developerMessage": err.Error(),
        "userMessage":      "unable to delete card",
    })
}

// POST /cards
//
// Params:
// id: a unique, positive integer that is the identifier of the assocoated deck
func CardPOST(db *sqlx.DB, ctx *gin.Context) {

    // parse request
    var (
        err         error
        jsonRequest CardPOSTRequest
    )

    err = ctx.BindJSON(&jsonRequest)
    if err != nil {

        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": err.Error(),
            "userMessage":      "bad JSON input",
        })
        ctx.Error(err)
        return
    }

    // fetch deck to verify it exists
    _, err = GetDeck(db, jsonRequest.Deck)
    switch {
    case err == ErrDeckNoSuchDeck:
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":           http.StatusNotFound,
            "developerMessage": err.Error(),
            "userMessage":      "canot find deck with given id",
        })
        ctx.Error(err)
        return
    case err != nil:
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve deck",
        })
        ctx.Error(err)
        return
    }

    // create card
    var newCardRow *CardRow

    newCardRow, err = CreateCard(db, &CardProps{
        Title:       jsonRequest.Title,
        Description: jsonRequest.Description,
        Front:       jsonRequest.Front,
        Back:        jsonRequest.Back,
        Deck:        jsonRequest.Deck,
    })
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to create new card",
        })
        ctx.Error(err)
    }

    var fetchedCardScore *CardScoreRow
    fetchedCardScore, err = GetCardScoreRecord(db, newCardRow.ID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card score record",
        })
        ctx.Error(err)
        return
    }

    var cardrow gin.H = CardRowToResponse(db, newCardRow)
    var cardscore gin.H = CardScoreToResponse(fetchedCardScore)

    ctx.JSON(http.StatusCreated, MergeResponses(
        &cardrow,
        &gin.H{"review": cardscore},
        &gin.H{"stashes": []uint{}},
    ))
}

// PATCH /cards/:id
//
// Input:
// title: non-empty string that shall be the new title of the deck
//
func CardPATCH(db *sqlx.DB, ctx *gin.Context) {

    var (
        err error
    )

    // parse id
    var cardIDString string = strings.ToLower(ctx.Param("id"))
    _cardID, err := strconv.ParseUint(cardIDString, 10, 32)
    var cardID uint = uint(_cardID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": err.Error(),
            "userMessage":      "given id is invalid",
        })
        ctx.Error(err)
        return
    }

    // parse request body
    var patch *StringMap = &StringMap{}
    err = ctx.BindJSON(patch)
    if err != nil {

        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": err.Error(),
            "userMessage":      "bad JSON input",
        })
        ctx.Error(err)
        return
    }
    if len(*patch) <= 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": "no JSON input",
            "userMessage":      "no JSON input",
        })
    }

    // TODO: validate patch
    // TODO: ensure title, if given, is non-empty string

    var (
        fetchedCardRow *CardRow = nil
    )

    // check requested card exists and fetch it

    fetchedCardRow, err = GetCard(db, cardID)
    switch {
    case err == ErrCardNoSuchCard:
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":           http.StatusNotFound,
            "developerMessage": err.Error(),
            "userMessage":      "unable to find card with given id",
        })
        ctx.Error(err)
        return
    case err != nil:
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to fetch card",
        })
        ctx.Error(err)
        return
    }

    var (
        deckID uint = fetchedCardRow.Deck
    )

    // case: moving card
    if _, hasDeckKey := (*patch)["deck"]; hasDeckKey == true {

        // check if new deck is valid

        // TODO: do this earlier for early bail
        // validate parent param
        var maybeDeckID uint
        maybeDeckID, err = (func() (uint, error) {
            switch _deckID := (*patch)["deck"].(type) {
            // note that according to docs: http://golang.org/pkg/encoding/json/#Unmarshal
            // JSON numbers are converted to float64
            case float64:
                if _deckID > 0 {
                    return uint(_deckID), nil
                }
            }
            return 0, errors.New("target deck is invalid")
        }())

        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "status":           http.StatusBadRequest,
                "developerMessage": err.Error(),
                "userMessage":      err.Error(),
            })
            ctx.Error(err)
            return
        }

        // cannot move card to the same parent
        if deckID == maybeDeckID {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "status":           http.StatusBadRequest,
                "developerMessage": "cannot move deck to the same parent",
                "userMessage":      "cannot move deck to the same parent",
            })
            return
        }
        deckID = maybeDeckID

        // validate target deck exists
        _, err = GetDeck(db, deckID)
        switch {
        case err == ErrDeckNoSuchDeck:
            ctx.JSON(http.StatusBadRequest, gin.H{
                "status":           http.StatusBadRequest,
                "developerMessage": err.Error(),
                "userMessage":      "given deck id is invalid",
            })
            ctx.Error(err)
            return
        case err != nil:
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "status":           http.StatusInternalServerError,
                "developerMessage": err.Error(),
                "userMessage":      "unable to fetch deck",
            })
            ctx.Error(err)
            return
        }
    }

    // generate SQL to patch card
    var (
        query string
        args  []interface{}
    )

    query, args, err = QueryApply(UPDATE_CARD_QUERY, &StringMap{"card_id": cardID}, patch)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to generate patch card SQL",
        })
        ctx.Error(err)
        return
    }

    var res sql.Result
    res, err = db.Exec(query, args...)
    if err != nil {
        // TODO: transaction rollback
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to patch card",
        })
        ctx.Error(err)
        return
    }

    // ensure card is patched
    num, err := res.RowsAffected()
    if err != nil {
        // TODO: transaction rollback
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to patch card",
        })
        ctx.Error(err)
        return
    }

    if num <= 0 {
        // TODO: transaction rollback
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":           http.StatusBadRequest,
            "developerMessage": "given JSON is invalid",
            "userMessage":      "given JSON is invalid",
        })
        return
    }

    fetchedCardRow, err = GetCard(db, cardID)
    switch {
    case err == ErrCardNoSuchCard:
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":           http.StatusNotFound,
            "developerMessage": err.Error(),
            "userMessage":      "cannot find card by id",
        })
        ctx.Error(err)
        return
    case err != nil:
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card",
        })
        ctx.Error(err)
        return
    }

    // fetch card score
    var fetchedCardScore *CardScoreRow
    fetchedCardScore, err = GetCardScoreRecord(db, cardID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card score record",
        })
        ctx.Error(err)
        return
    }

    var fetchedStashes []uint
    fetchedStashes, err = StashesByCard(db, cardID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":           http.StatusInternalServerError,
            "developerMessage": err.Error(),
            "userMessage":      "unable to retrieve card stashes",
        })
        ctx.Error(err)
        return
    }

    var cardrow gin.H = CardRowToResponse(db, fetchedCardRow)
    var cardscore gin.H = CardScoreToResponse(fetchedCardScore)

    ctx.JSON(http.StatusOK, MergeResponses(
        &cardrow,
        &gin.H{"review": cardscore},
        &gin.H{"stashes": fetchedStashes},
    ))
}

/* helpers */

func CardResponse(overrides *gin.H) gin.H {
    defaultResponse := &gin.H{
        "id":          0,  // required
        "title":       "", // required
        "description": "",
        "front":       "",
        "back":        "",
        "deck":        0,  // required
        "created_at":  0,
        "updated_at":  0,
        "deck_path":   []uint{},
    }

    return MergeResponse(defaultResponse, overrides)
}

func CardRowToResponse(db *sqlx.DB, cardrow *CardRow) gin.H {

    // construct deck_path
    // var err error
    var deck_path []uint
    deck_path, _ = GetDeckAncestors(db, cardrow.Deck)
    // swallow error
    // TODO: error handling

    deck_path = append(deck_path, cardrow.Deck)

    return CardResponse(&gin.H{
        "id":          cardrow.ID,
        "title":       cardrow.Title,
        "description": cardrow.Description,
        "front":       cardrow.Front,
        "back":        cardrow.Back,
        "deck":        cardrow.Deck,
        "created_at":  cardrow.CreatedAt,
        "updated_at":  cardrow.UpdatedAt,
        "deck_path":   deck_path,
    })
}

func ValidateCardProps(props *CardProps) error {

    if len(props.Title) <= 0 {
        return errors.New("Title must be non-empty string")
    }

    if props.Deck <= 0 {
        return errors.New("Deck id must be positive non-zero integer")
    }

    return nil
}

func GetCard(db *sqlx.DB, cardID uint) (*CardRow, error) {

    var (
        err   error
        query string
        args  []interface{}
    )

    query, args, err = QueryApply(FETCH_CARD_QUERY, &StringMap{"card_id": cardID})
    if err != nil {
        return nil, err
    }

    var fetchedCard *CardRow = &CardRow{}

    err = db.QueryRowx(query, args...).StructScan(fetchedCard)

    switch {
    case err == sql.ErrNoRows:
        return nil, ErrCardNoSuchCard
    case err != nil:
        return nil, err
    default:
        return fetchedCard, nil
    }
}

func CreateCard(db *sqlx.DB, props *CardProps) (*CardRow, error) {

    var err error

    // validate card props
    err = ValidateCardProps(props)
    if err != nil {
        return nil, err
    }

    var (
        res   sql.Result
        query string
        args  []interface{}
    )

    query, args, err = QueryApply(CREATE_NEW_CARD_QUERY,
        &StringMap{
            "title":       props.Title,
            "description": props.Description,
            "front":       props.Front,
            "back":        props.Back,
            "deck":        props.Deck,
        })
    if err != nil {
        return nil, err
    }

    // insert new card
    res, err = db.Exec(query, args...)
    if err != nil {
        return nil, err
    }

    insertID, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    return GetCard(db, uint(insertID))
}

func CountCardsByDeck(db *sqlx.DB, deckID uint) (uint, error) {

    var (
        err   error
        query string
        args  []interface{}
    )

    query, args, err = QueryApply(COUNT_CARDS_BY_DECK_QUERY, &StringMap{
        "deck_id": deckID,
    })
    if err != nil {
        return 0, err
    }

    var count uint
    err = db.QueryRowx(query, args...).Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func CardsByDeck(db *sqlx.DB, queryTransform PipeInput, deckID uint, page uint, per_page uint) (*([]CardRow), error) {

    var (
        err   error
        query string
        args  []interface{}
    )

    // invariant: page >= 1

    var offset uint = (page - 1) * per_page

    var count uint
    count, err = CountCardsByDeck(db, deckID)
    if err != nil {
        return nil, err
    }

    if count <= 0 {
        return nil, ErrCardNoCardsByDeck
    }

    if offset >= count {
        return nil, ErrCardPageOutOfBounds
    }

    query, args, err = QueryApply(queryTransform, &StringMap{
        "deck_id":  deckID,
        "per_page": per_page,
        "offset":   offset,
    })
    if err != nil {
        return nil, err
    }

    var cards []CardRow = make([]CardRow, 0, per_page)
    err = db.Select(&cards, query, args...)
    if err != nil {
        return nil, err
    }

    if len(cards) <= 0 {
        return nil, ErrCardNoCardsByDeck
    }

    return &cards, nil
}

func DeleteCard(db *sqlx.DB, cardID uint) error {

    var (
        err   error
        query string
        args  []interface{}
    )

    query, args, err = QueryApply(DELETE_CARD_QUERY, &StringMap{"card_id": cardID})
    if err != nil {
        return err
    }

    _, err = db.Exec(query, args...)
    if err != nil {
        return err
    }

    return nil
}
