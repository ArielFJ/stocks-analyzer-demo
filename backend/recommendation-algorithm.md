# Stock Recommendation Algorithm Flow

## Main Algorithm Flow

```mermaid
flowchart TB
    A[API Request: /stocks/recommendations] --> B[StockService.GetRecommendations]
    B --> C[Get All Stocks with Analysis Data]
    C --> D[RecommendationEngine.AnalyzeStocks]
    
    D --> E{For Each Stock}
    E --> F{Has Analysis Data?}
    F -->|No| G[Skip Stock]
    F -->|Yes| H[calculateScore - Start with Base Score 50.0]
    
    H --> SCORE[Score Calculation Process]
    SCORE --> CONF[Confidence & Reason Generation]
    CONF --> REC[Create StockRecommendation]
    
    REC --> V{More Stocks?}
    V -->|Yes| E
    V -->|No| W[Sort by Score DESC]
    
    W --> X[Take Top 10]
    X --> Y[Return Recommendations]
    
    G --> V
    
    style A fill:#e1f5fe
    style Y fill:#c8e6c9
    style H fill:#fff3e0
    style SCORE fill:#ffecb3
    style CONF fill:#f3e5f5
```

## Score Calculation Detail

```mermaid
flowchart TB
    START[Base Score: 50.0] --> RATING[Rating Analysis]
    START --> TARGET[Price Target Analysis]
    START --> ACTION[Action Type Analysis]
    START --> COVERAGE[Coverage Depth Analysis]
    
    RATING --> R1{Rating Change?}
    R1 -->|Upgrade| R2[+15 bonus]
    R1 -->|Downgrade| R3[-10 penalty]
    R1 -->|No Change| R4[Add rating score only]
    
    TARGET --> T1{Target Change?}
    T1 -->|Over 10% increase| T2[+20 points]
    T1 -->|5-10% increase| T3[+10 points]
    T1 -->|Over 10% decrease| T4[-15 points]
    T1 -->|5-10% decrease| T5[-8 points]
    
    ACTION --> A1{Action Type?}
    A1 -->|Initiated| A2[+10 points]
    A1 -->|Raised| A3[+12 points]
    A1 -->|Lowered| A4[-8 points]
    A1 -->|Maintained| A5[+5 points]
    
    COVERAGE --> C1{≥3 Recent Analyses?}
    C1 -->|Yes| C2[+5 points]
    COVERAGE --> C3{≥2 Positive Ratings?}
    C3 -->|Yes| C4[+8 points]
    
    R2 --> FINAL[Calculate Final Score]
    R3 --> FINAL
    R4 --> FINAL
    T2 --> FINAL
    T3 --> FINAL
    T4 --> FINAL
    T5 --> FINAL
    A2 --> FINAL
    A3 --> FINAL
    A4 --> FINAL
    A5 --> FINAL
    C2 --> FINAL
    C4 --> FINAL
    
    style START fill:#e3f2fd
    style FINAL fill:#c8e6c9
```

## Confidence Assignment

```mermaid
flowchart LR
    SCORE[Final Score] --> CHECK1{Score ≥ 75?}
    CHECK1 -->|Yes| HIGH[High Confidence]
    CHECK1 -->|No| CHECK2{Score ≥ 60?}
    CHECK2 -->|Yes| MEDIUM[Medium Confidence]
    CHECK2 -->|No| LOW[Low Confidence]
    
    style HIGH fill:#4caf50
    style MEDIUM fill:#ff9800
    style LOW fill:#f44336
```

## Reason Generation Process

```mermaid
flowchart TB
    START[Generate Reason] --> CHECK_RATING[Check Buy/Outperform Rating]
    START --> CHECK_TARGET[Check Price Target Increase]
    START --> CHECK_COVERAGE[Check New Coverage]
    START --> CHECK_MULTIPLE[Check Multiple Updates]
    
    CHECK_RATING --> REASON1[Buy rating from Brokerage]
    CHECK_TARGET --> REASON2[Price target raised by X%]
    CHECK_COVERAGE --> REASON3[New analyst coverage]
    CHECK_MULTIPLE --> REASON4[Multiple recent updates]
    
    REASON1 --> COMBINE[Combine Reasons]
    REASON2 --> COMBINE
    REASON3 --> COMBINE
    REASON4 --> COMBINE
    
    COMBINE --> FINAL_REASON[Final Reason String]
    
    style START fill:#e8f5e8
    style FINAL_REASON fill:#c8e6c9
```

## Rating Score Mapping

```mermaid
graph LR
    A[Rating Input] --> B{Parse Rating}
    B -->|Strong Buy/Buy| C[80 points]
    B -->|Outperform/Overweight| D[70 points]
    B -->|Hold/Neutral| E[50 points]
    B -->|Underperform/Underweight| F[30 points]
    B -->|Sell/Strong Sell| G[10 points]
    B -->|Unknown| H[50 points default]
    
    style C fill:#4caf50
    style D fill:#8bc34a
    style E fill:#ffeb3b
    style F fill:#ff9800
    style G fill:#f44336
    style H fill:#9e9e9e
```

## Algorithm Summary

The recommendation algorithm processes stocks through a multi-factor scoring system:

1. **Base Score**: Every stock starts with 50 points
2. **Rating Analysis**: Adds/subtracts based on analyst ratings and changes
3. **Price Target Analysis**: Rewards target increases, penalizes decreases
4. **Action Analysis**: Considers the type of analyst action taken
5. **Coverage Analysis**: Rewards multiple analyses and positive sentiment
6. **Confidence Assignment**: Categorizes based on final score
7. **Reason Generation**: Creates human-readable explanations
8. **Ranking**: Sorts by score and returns top 10 recommendations

The algorithm emphasizes recent positive analyst actions and upgrades, making it effective for identifying stocks with improving market sentiment.